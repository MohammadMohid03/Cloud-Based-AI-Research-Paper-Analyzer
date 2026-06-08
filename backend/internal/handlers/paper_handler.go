package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/middleware"
	"github.com/research-paper-analyzer/backend/internal/models"
	"github.com/research-paper-analyzer/backend/internal/services/paper"
	"github.com/research-paper-analyzer/backend/internal/services/s3"
	"gorm.io/gorm"
)

// PaperHandler handles research paper related HTTP requests.
type PaperHandler struct {
	db           *gorm.DB
	paperService *paper.Service
	s3Service    *s3.Service
}

// NewPaperHandler creates a new PaperHandler instance.
func NewPaperHandler(db *gorm.DB, paperService *paper.Service, s3Service *s3.Service) *PaperHandler {
	return &PaperHandler{
		db:           db,
		paperService: paperService,
		s3Service:    s3Service,
	}
}

// ListPapers returns a list of papers belonging to the authenticated user.
// GET /api/papers
func (h *PaperHandler) ListPapers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	papers, err := h.paperService.GetUserPapers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"papers": papers})
}

// GetPaper returns detailed information for a specific paper, including analysis results.
// GET /api/papers/:id
func (h *PaperHandler) GetPaper(c *gin.Context) {
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

	p, err := h.paperService.GetPaperWithAnalysis(uint(paperID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// UploadPaper handles the PDF file upload and stores the metadata.
// POST /api/papers/upload
func (h *PaperHandler) UploadPaper(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Retrieve the file from form-data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Verify it's a PDF
	if header.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are supported"})
		return
	}

	// Retrieve title (or use filename as fallback)
	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	// Upload to configured storage
	fileURL, s3Key, err := h.s3Service.UploadFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Create a new paper record in DB
	newPaper := models.Paper{
		UserID:  userID,
		Title:   title,
		FileURL: fileURL,
		S3Key:   s3Key,
		Status:  models.StatusPending,
	}

	if err := h.db.Create(&newPaper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save paper metadata"})
		return
	}

	c.JSON(http.StatusCreated, newPaper)
}

// DeletePaper handles deleting a paper from the database.
// DELETE /api/papers/:id
func (h *PaperHandler) DeletePaper(c *gin.Context) {
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

	// Verify ownership and delete
	var paper models.Paper
	result := h.db.Where("id = ? AND user_id = ?", paperID, userID).First(&paper)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	if err := h.db.Delete(&paper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete paper"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paper deleted successfully"})
}
