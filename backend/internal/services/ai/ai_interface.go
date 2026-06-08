// Package ai defines the interface for AI services and provides
// implementations for both mock testing and AWS Bedrock production use.
package ai

import (
	"github.com/research-paper-analyzer/backend/internal/models"
)

// AnalysisResult holds the structured output from AI analysis of a paper.
type AnalysisResult struct {
	Summary     string // Comprehensive summary of the paper
	KeyFindings string // Bullet-point list of key findings
	Methodology string // Description of the research methodology
	Limitations string // Identified limitations of the research
	FutureScope string // Future research directions
	Keywords    string // Comma-separated relevant keywords
}

// AIService defines the interface that all AI providers must implement.
// This abstraction allows easy switching between mock and production AI services.
type AIService interface {
	// GenerateSummary analyzes paper text and returns structured analysis results.
	GenerateSummary(text string) (*AnalysisResult, error)

	// GenerateQuiz creates multiple-choice questions based on the paper text.
	GenerateQuiz(text string, numQuestions int) ([]models.QuizQuestionAI, error)

	// ChatWithContext answers a question using paper chunks as context,
	// with awareness of previous chat history.
	ChatWithContext(question string, context string, chatHistory []models.ChatMessage) (string, error)
}
