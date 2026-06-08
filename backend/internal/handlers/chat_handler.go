package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/middleware"
	"github.com/research-paper-analyzer/backend/internal/services/chat"
)

// ChatHandler handles chat requests about papers.
type ChatHandler struct {
	chatService *chat.Service
}

// NewChatHandler creates a new ChatHandler instance.
func NewChatHandler(chatService *chat.Service) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// ChatWithPaper handles user questions about a paper.
// POST /api/papers/:id/chat
func (h *ChatHandler) ChatWithPaper(c *gin.Context) {
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

	// Bind request JSON
	var req struct {
		Question string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question is required"})
		return
	}

	response, err := h.chatService.AskQuestion(uint(paperID), userID, req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
