package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/middleware"
	"github.com/research-paper-analyzer/backend/internal/services/quiz"
)

// QuizHandler handles quiz question retrieval and generation.
type QuizHandler struct {
	quizService *quiz.Service
}

// NewQuizHandler creates a new QuizHandler instance.
func NewQuizHandler(quizService *quiz.Service) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

// GetQuiz retrieves or generates quiz questions for a specific paper.
// GET /api/papers/:id/quiz
func (h *QuizHandler) GetQuiz(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	paperIDStr := c.Param("id")
	paperID, err := strconv.ParseUint(paperIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid paper ID"})
		return
	}

	questions, err := h.quizService.GetOrGenerateQuiz(uint(paperID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}
