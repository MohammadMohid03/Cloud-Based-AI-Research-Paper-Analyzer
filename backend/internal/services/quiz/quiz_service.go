// Package quiz provides quiz generation and retrieval services.
// It uses the AI service to generate questions and stores them in the database.
package quiz

import (
	"fmt"
	"log"

	"github.com/research-paper-analyzer/backend/internal/models"
	"github.com/research-paper-analyzer/backend/internal/services/ai"
	"gorm.io/gorm"
)

// Service handles quiz question generation and retrieval.
type Service struct {
	db        *gorm.DB
	aiService ai.AIService
}

// NewService creates a new quiz service instance.
func NewService(db *gorm.DB, aiService ai.AIService) *Service {
	return &Service{
		db:        db,
		aiService: aiService,
	}
}

// GetOrGenerateQuiz retrieves existing quiz questions for a paper,
// or generates new ones if none exist.
func (s *Service) GetOrGenerateQuiz(paperID uint, userID uint) ([]models.QuizQuestion, error) {
	// Verify the paper belongs to the user
	var paper models.Paper
	if err := s.db.Where("id = ? AND user_id = ?", paperID, userID).First(&paper).Error; err != nil {
		return nil, fmt.Errorf("paper not found or access denied: %w", err)
	}

	// Check if quiz questions already exist
	var existingQuestions []models.QuizQuestion
	s.db.Where("paper_id = ?", paperID).Find(&existingQuestions)

	if len(existingQuestions) > 0 {
		return existingQuestions, nil
	}

	// Generate new quiz questions if the paper has been analyzed
	if paper.Status != models.StatusCompleted {
		return nil, fmt.Errorf("paper must be analyzed before generating quiz questions (current status: %s)", paper.Status)
	}

	// Get the paper text for quiz generation
	text := paper.RawText
	if text == "" {
		return nil, fmt.Errorf("no text available for quiz generation")
	}

	log.Printf("🎓 Generating quiz questions for paper ID %d", paperID)

	// Call AI service to generate quiz questions
	aiQuestions, err := s.aiService.GenerateQuiz(text, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to generate quiz: %w", err)
	}

	// Store generated questions in the database
	var savedQuestions []models.QuizQuestion
	for _, q := range aiQuestions {
		question := models.QuizQuestion{
			PaperID:       paperID,
			Question:      q.Question,
			OptionA:       q.OptionA,
			OptionB:       q.OptionB,
			OptionC:       q.OptionC,
			OptionD:       q.OptionD,
			CorrectAnswer: q.CorrectAnswer,
			Explanation:   q.Explanation,
		}

		if err := s.db.Create(&question).Error; err != nil {
			log.Printf("⚠️  Warning: failed to save quiz question: %v", err)
			continue
		}
		savedQuestions = append(savedQuestions, question)
	}

	log.Printf("✅ Generated and saved %d quiz questions for paper ID %d", len(savedQuestions), paperID)
	return savedQuestions, nil
}
