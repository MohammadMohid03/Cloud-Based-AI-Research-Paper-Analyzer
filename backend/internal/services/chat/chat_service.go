// Package chat provides the chat-with-paper service.
// It retrieves relevant paper chunks, builds context, and uses the AI service
// to generate answers to user questions about the paper.
package chat

import (
	"fmt"
	"log"
	"strings"

	"github.com/research-paper-analyzer/backend/internal/models"
	"github.com/research-paper-analyzer/backend/internal/services/ai"
	"gorm.io/gorm"
)

// Service handles chat interactions with papers.
type Service struct {
	db        *gorm.DB
	aiService ai.AIService
}

// NewService creates a new chat service instance.
func NewService(db *gorm.DB, aiService ai.AIService) *Service {
	return &Service{
		db:        db,
		aiService: aiService,
	}
}

// AskQuestion processes a user's question about a paper:
// 1. Finds relevant chunks using keyword matching
// 2. Retrieves recent chat history for context
// 3. Sends the question + context to the AI service
// 4. Stores and returns the response
func (s *Service) AskQuestion(paperID uint, userID uint, question string) (*models.ChatResponse, error) {
	// Verify the paper belongs to the user
	var paper models.Paper
	if err := s.db.Where("id = ? AND user_id = ?", paperID, userID).First(&paper).Error; err != nil {
		return nil, fmt.Errorf("paper not found or access denied: %w", err)
	}

	// Step 1: Find relevant chunks using keyword matching
	relevantContext, err := s.findRelevantChunks(paperID, question)
	if err != nil {
		log.Printf("⚠️  Warning: failed to find relevant chunks: %v", err)
		// Fall back to the first few chunks or raw text
		relevantContext = truncateText(paper.RawText, 3000)
	}

	// Step 2: Retrieve recent chat history (last 5 exchanges)
	chatHistory := s.getChatHistory(paperID, userID, 5)

	// Step 3: Send to AI service
	answer, err := s.aiService.ChatWithContext(question, relevantContext, chatHistory)
	if err != nil {
		return nil, fmt.Errorf("AI chat failed: %w", err)
	}

	// Step 4: Store the chat exchange
	chatEntry := models.ChatHistory{
		PaperID:  paperID,
		UserID:   userID,
		Question: question,
		Answer:   answer,
	}

	if err := s.db.Create(&chatEntry).Error; err != nil {
		log.Printf("⚠️  Warning: failed to save chat history: %v", err)
	}

	return &models.ChatResponse{
		Question:  question,
		Answer:    answer,
		CreatedAt: chatEntry.CreatedAt,
	}, nil
}

// findRelevantChunks uses simple keyword matching to find the most relevant
// text chunks for answering a question. Each chunk is scored based on
// how many question keywords it contains.
func (s *Service) findRelevantChunks(paperID uint, question string) (string, error) {
	// Get all chunks for this paper
	var chunks []models.PaperChunk
	if err := s.db.Where("paper_id = ?", paperID).Order("chunk_index ASC").Find(&chunks).Error; err != nil {
		return "", fmt.Errorf("failed to fetch chunks: %w", err)
	}

	if len(chunks) == 0 {
		return "", fmt.Errorf("no chunks found for paper")
	}

	// Extract keywords from the question (simple approach: split and filter stop words)
	keywords := extractKeywords(question)

	// Score each chunk based on keyword matches
	type scoredChunk struct {
		chunk models.PaperChunk
		score int
	}

	var scored []scoredChunk
	for _, chunk := range chunks {
		score := 0
		chunkLower := strings.ToLower(chunk.ChunkText)
		for _, keyword := range keywords {
			// Count occurrences of each keyword in the chunk
			score += strings.Count(chunkLower, strings.ToLower(keyword))
		}
		scored = append(scored, scoredChunk{chunk: chunk, score: score})
	}

	// Sort by score (descending) - simple bubble sort for clarity
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Take top 3 most relevant chunks (or all if fewer)
	maxChunks := 3
	if len(scored) < maxChunks {
		maxChunks = len(scored)
	}

	var contextParts []string
	for i := 0; i < maxChunks; i++ {
		if scored[i].score > 0 || i == 0 {
			// Always include at least one chunk, and only include others if they match keywords
			contextParts = append(contextParts, scored[i].chunk.ChunkText)
		}
	}

	return strings.Join(contextParts, "\n\n---\n\n"), nil
}

// getChatHistory retrieves recent chat exchanges for context continuity.
func (s *Service) getChatHistory(paperID uint, userID uint, limit int) []models.ChatMessage {
	var history []models.ChatHistory
	s.db.Where("paper_id = ? AND user_id = ?", paperID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history)

	// Convert to ChatMessage format and reverse (so oldest first)
	messages := make([]models.ChatMessage, 0, len(history)*2)
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, models.ChatMessage{
			Role:    "user",
			Content: history[i].Question,
		})
		messages = append(messages, models.ChatMessage{
			Role:    "assistant",
			Content: history[i].Answer,
		})
	}

	return messages
}

// extractKeywords splits a question into meaningful keywords by removing
// common English stop words.
func extractKeywords(text string) []string {
	// Common English stop words to filter out
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "shall": true, "can": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "at": true, "by": true, "from": true, "as": true,
		"into": true, "about": true, "between": true, "through": true,
		"and": true, "but": true, "or": true, "nor": true, "not": true,
		"so": true, "yet": true, "both": true, "either": true,
		"i": true, "me": true, "my": true, "we": true, "our": true,
		"you": true, "your": true, "he": true, "she": true, "it": true,
		"they": true, "them": true, "their": true, "its": true,
		"this": true, "that": true, "these": true, "those": true,
		"what": true, "which": true, "who": true, "whom": true, "how": true,
		"when": true, "where": true, "why": true,
		"if": true, "then": true, "else": true, "than": true,
		"very": true, "just": true, "also": true, "there": true,
	}

	words := strings.Fields(strings.ToLower(text))
	var keywords []string

	for _, word := range words {
		// Remove punctuation from word edges
		word = strings.Trim(word, ".,!?;:'\"()[]{}—-")
		// Keep words that are meaningful (not stop words, minimum length)
		if len(word) >= 3 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// truncateText truncates text to a maximum number of characters,
// breaking at word boundaries.
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	// Find the last space to avoid breaking mid-word
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}
