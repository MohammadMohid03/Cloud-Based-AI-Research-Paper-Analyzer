package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/middleware"
	"github.com/research-paper-analyzer/backend/internal/models"
	"github.com/research-paper-analyzer/backend/internal/services/paper"
	"gorm.io/gorm"
)

// AnalysisHandler handles analysis triggers and report generation.
type AnalysisHandler struct {
	db           *gorm.DB
	paperService *paper.Service
}

// NewAnalysisHandler creates a new AnalysisHandler instance.
func NewAnalysisHandler(db *gorm.DB, paperService *paper.Service) *AnalysisHandler {
	return &AnalysisHandler{
		db:           db,
		paperService: paperService,
	}
}

// AnalyzePaper triggers the asynchronous PDF extraction and AI analysis pipeline.
// POST /api/papers/:id/analyze
func (h *AnalysisHandler) AnalyzePaper(c *gin.Context) {
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

	// Verify ownership first
	var p models.Paper
	if err := h.db.Where("id = ? AND user_id = ?", paperID, userID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	// Trigger analysis in a background goroutine to avoid timing out the HTTP request
	go func(id uint) {
		if err := h.paperService.AnalyzePaper(id); err != nil {
			// GORM DB writes inside AnalyzePaper will update status to Failed
			println("Analysis error for paper", id, ":", err.Error())
		}
	}(uint(paperID))

	c.JSON(http.StatusOK, gin.H{
		"message": "Analysis started successfully. Please poll/refresh paper status.",
		"status":  models.StatusProcessing,
	})
}

// GetReport generates and returns the analysis results as a JSON download/response.
// GET /api/papers/:id/report
func (h *AnalysisHandler) GetReport(c *gin.Context) {
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

	// Fetch paper with preloaded Analysis
	var p models.Paper
	err = h.db.Preload("Analysis").Where("id = ? AND user_id = ?", paperID, userID).First(&p).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	if p.Status != models.StatusCompleted || p.Analysis == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Analysis is not completed yet"})
		return
	}

	// Set headers to trigger file download in browser
	c.Header("Content-Disposition", "attachment; filename=analysis_report_"+strconv.Itoa(int(p.ID))+".json")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, p.Analysis)
}
