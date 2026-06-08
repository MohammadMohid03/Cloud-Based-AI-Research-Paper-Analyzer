package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/research-paper-analyzer/backend/internal/config"
	"github.com/research-paper-analyzer/backend/internal/database"
	"github.com/research-paper-analyzer/backend/internal/handlers"
	"github.com/research-paper-analyzer/backend/internal/routes"
	"github.com/research-paper-analyzer/backend/internal/services/ai"
	"github.com/research-paper-analyzer/backend/internal/services/auth"
	"github.com/research-paper-analyzer/backend/internal/services/chat"
	"github.com/research-paper-analyzer/backend/internal/services/paper"
	"github.com/research-paper-analyzer/backend/internal/services/pdf"
	"github.com/research-paper-analyzer/backend/internal/services/quiz"
	"github.com/research-paper-analyzer/backend/internal/services/s3"
)

func main() {
	log.Println("🚀 Starting AI-Powered Research Paper Analyzer Backend...")

	// 1. Load configuration from environment variables
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// 2. Initialize PostgreSQL database connection
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("❌ Database initialization failed: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("⚠️ Error closing database connection: %v", err)
		}
	}()

	// 3. Run GORM auto-migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	// 4. Seed database with development/demo data if enabled
	if cfg.SeedData {
		if err := database.SeedData(); err != nil {
			log.Printf("⚠️ Database seeding failed: %v", err)
		}
	}

	// 5. Initialize S3 or Local storage service
	s3Svc, err := s3.NewService(cfg)
	if err != nil {
		log.Fatalf("❌ Storage service initialization failed: %v", err)
	}

	// 6. Initialize PDF text extraction service
	pdfSvc := pdf.NewService()

	// 7. Initialize AI service (Mock or Bedrock)
	var aiSvc ai.AIService
	if cfg.AIProvider == "bedrock" {
		var err error
		aiSvc, err = ai.NewBedrockAIService(cfg)
		if err != nil {
			log.Printf("⚠️ Failed to initialize Amazon Bedrock Service: %v. Falling back to Mock AI.", err)
			aiSvc = ai.NewMockAIService()
		} else {
			log.Println("✅ Amazon Bedrock AI Service initialized successfully")
		}
	} else {
		aiSvc = ai.NewMockAIService()
		log.Println("✅ Mock AI Service initialized successfully (local testing mode)")
	}

	// 8. Initialize core business services
	db := database.GetDB()
	authSvc := auth.NewService(db, cfg)
	paperSvc := paper.NewService(db, aiSvc, pdfSvc)
	chatSvc := chat.NewService(db, aiSvc)
	quizSvc := quiz.NewService(db, aiSvc)

	// 9. Initialize HTTP Handlers
	authH := handlers.NewAuthHandler(authSvc)
	paperH := handlers.NewPaperHandler(db, paperSvc, s3Svc)
	analysisH := handlers.NewAnalysisHandler(db, paperSvc)
	chatH := handlers.NewChatHandler(chatSvc)
	quizH := handlers.NewQuizHandler(quizSvc)

	// 10. Initialize Gin router and setup endpoints
	r := gin.New()
	routes.SetupRoutes(r, cfg, authH, paperH, analysisH, chatH, quizH)

	// Fallback route for 404 page
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
	})

	// 11. Run the HTTP server
	addr := ":" + cfg.Port
	log.Printf("📡 Server listening on %s in %s mode", addr, cfg.GinMode)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
