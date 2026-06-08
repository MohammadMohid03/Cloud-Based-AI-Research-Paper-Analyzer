package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/config"
	"github.com/research-paper-analyzer/backend/internal/handlers"
	"github.com/research-paper-analyzer/backend/internal/middleware"
)

// SetupRoutes registers all API routes and middlewares in the Gin engine.
func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authH *handlers.AuthHandler,
	paperH *handlers.PaperHandler,
	analysisH *handlers.AnalysisHandler,
	chatH *handlers.ChatHandler,
	quizH *handlers.QuizHandler,
) {
	// Enable CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Enable Logger and Recovery middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())

	// API base group
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})

		// Auth group (Public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}

		// Protected routes group (Requires JWT Auth)
		protected := api.Group("/papers")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("", paperH.ListPapers)
			protected.POST("/upload", paperH.UploadPaper)
			protected.GET("/:id", paperH.GetPaper)
			protected.DELETE("/:id", paperH.DeletePaper)
			protected.POST("/:id/analyze", analysisH.AnalyzePaper)
			protected.POST("/:id/chat", chatH.ChatWithPaper)
			protected.GET("/:id/quiz", quizH.GetQuiz)
			protected.GET("/:id/report", analysisH.GetReport)
		}
	}
}
