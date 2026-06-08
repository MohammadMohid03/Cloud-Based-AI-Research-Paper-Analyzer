// Package paper provides the paper analysis orchestration service.
// It coordinates PDF extraction, chunking, AI analysis, and database storage.
package paper

import (
	"fmt"
	"log"

	"github.com/research-paper-analyzer/backend/internal/models"
	"github.com/research-paper-analyzer/backend/internal/services/ai"
	"github.com/research-paper-analyzer/backend/internal/services/pdf"
	"gorm.io/gorm"
)

// Service orchestrates the full paper analysis pipeline.
type Service struct {
	db         *gorm.DB
	aiService  ai.AIService
	pdfService *pdf.Service
}

// NewService creates a new paper service instance.
func NewService(db *gorm.DB, aiService ai.AIService, pdfService *pdf.Service) *Service {
	return &Service{
		db:         db,
		aiService:  aiService,
		pdfService: pdfService,
	}
}

// AnalyzePaper runs the full analysis pipeline for a paper:
// 1. Extracts text from the PDF
// 2. Splits text into chunks and stores them
// 3. Calls the AI service for structured analysis
// 4. Stores analysis results in the database
// 5. Updates the paper status
func (s *Service) AnalyzePaper(paperID uint) error {
	// Fetch the paper from the database
	var paper models.Paper
	if err := s.db.First(&paper, paperID).Error; err != nil {
		return fmt.Errorf("paper not found: %w", err)
	}

	// Update status to processing
	s.db.Model(&paper).Update("status", models.StatusProcessing)
	log.Printf("🔄 Starting analysis for paper ID %d: %s", paper.ID, paper.Title)

	// Step 1: Extract text from PDF (if not already extracted)
	if paper.RawText == "" {
		log.Printf("  📄 Extracting text from PDF: %s", paper.FileURL)
		extractedText, err := s.pdfService.ExtractText(paper.FileURL)
		if err != nil {
			s.db.Model(&paper).Update("status", models.StatusFailed)
			return fmt.Errorf("failed to extract text: %w", err)
		}
		paper.RawText = extractedText
		s.db.Model(&paper).Update("raw_text", extractedText)
		log.Printf("  ✅ Text extracted (%d characters)", len(extractedText))
	}

	// Step 2: Split into chunks and store them
	log.Println("  📝 Splitting text into chunks...")
	chunks := s.pdfService.SplitIntoChunks(paper.RawText, 500, 50)

	// Delete any existing chunks for this paper (in case of re-analysis)
	s.db.Where("paper_id = ?", paper.ID).Delete(&models.PaperChunk{})

	for i, chunkText := range chunks {
		chunk := models.PaperChunk{
			PaperID:    paper.ID,
			ChunkText:  chunkText,
			ChunkIndex: i,
		}
		if err := s.db.Create(&chunk).Error; err != nil {
			log.Printf("  ⚠️  Warning: failed to save chunk %d: %v", i, err)
		}
	}
	log.Printf("  ✅ Created %d chunks", len(chunks))

	// Step 3: Call AI service for analysis
	log.Println("  🤖 Running AI analysis...")
	analysisResult, err := s.aiService.GenerateSummary(paper.RawText)
	if err != nil {
		s.db.Model(&paper).Update("status", models.StatusFailed)
		return fmt.Errorf("AI analysis failed: %w", err)
	}

	// Step 4: Store analysis results
	// Delete any existing analysis (for re-analysis)
	s.db.Where("paper_id = ?", paper.ID).Delete(&models.PaperAnalysis{})

	analysis := models.PaperAnalysis{
		PaperID:     paper.ID,
		Summary:     analysisResult.Summary,
		KeyFindings: analysisResult.KeyFindings,
		Methodology: analysisResult.Methodology,
		Limitations: analysisResult.Limitations,
		FutureScope: analysisResult.FutureScope,
		Keywords:    analysisResult.Keywords,
	}

	if err := s.db.Create(&analysis).Error; err != nil {
		s.db.Model(&paper).Update("status", models.StatusFailed)
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	// Step 5: Update paper status to completed
	s.db.Model(&paper).Update("status", models.StatusCompleted)
	log.Printf("✅ Analysis completed for paper ID %d", paper.ID)

	return nil
}

// GetPaperWithAnalysis fetches a paper with all its related data.
func (s *Service) GetPaperWithAnalysis(paperID uint, userID uint) (*models.Paper, error) {
	var paper models.Paper
	err := s.db.
		Preload("Analysis").
		Preload("Chunks", func(db *gorm.DB) *gorm.DB {
			return db.Order("chunk_index ASC")
		}).
		Where("id = ? AND user_id = ?", paperID, userID).
		First(&paper).Error

	if err != nil {
		return nil, fmt.Errorf("paper not found: %w", err)
	}

	return &paper, nil
}

// GetUserPapers returns all papers belonging to a user (lightweight, without full text).
func (s *Service) GetUserPapers(userID uint) ([]models.PaperResponse, error) {
	var papers []models.Paper
	err := s.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&papers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch papers: %w", err)
	}

	// Convert to lightweight response
	responses := make([]models.PaperResponse, len(papers))
	for i, p := range papers {
		responses[i] = p.ToPaperResponse()
	}

	return responses, nil
}
