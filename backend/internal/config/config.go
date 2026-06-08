// Package config handles loading and providing access to application configuration.
// It reads environment variables (optionally from a .env file) and exposes them
// through a strongly-typed Config struct.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	// Server settings
	Port    string // HTTP server port (default: "8080")
	GinMode string // Gin framework mode: "debug" or "release"

	// Database settings
	DBHost     string // PostgreSQL host
	DBPort     string // PostgreSQL port
	DBUser     string // PostgreSQL username
	DBPassword string // PostgreSQL password
	DBName     string // PostgreSQL database name
	DBSSLMode  string // PostgreSQL SSL mode (disable, require, etc.)

	// JWT settings
	JWTSecret      string // Secret key for signing JWT tokens
	JWTExpiryHours int    // Token expiration time in hours

	// AI provider settings
	AIProvider    string // "mock" or "bedrock"
	BedrockModel  string // AWS Bedrock model ID
	AWSRegion     string // AWS region
	AWSAccessKey  string // AWS access key ID (optional, can use IAM roles)
	AWSSecretKey  string // AWS secret access key

	// Storage settings
	StorageProvider string // "local" or "s3"
	S3BucketName    string // S3 bucket name for file uploads
	UploadDir       string // Local directory for file uploads

	// CORS settings
	CORSAllowedOrigins []string // List of allowed origins for CORS

	// Seed data
	SeedData bool // Whether to seed demo data on startup
}

// Load reads configuration from environment variables and returns a Config struct.
// It attempts to load a .env file first (useful for local development),
// but won't fail if the file doesn't exist (e.g., in production with real env vars).
func Load() *Config {
	// Try to load .env file - don't fail if it doesn't exist
	// In production, environment variables are set directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	// Parse JWT expiry hours with a sensible default
	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	if err != nil {
		jwtExpiry = 72 // Default to 72 hours if parsing fails
	}

	// Parse seed data flag
	seedData := strings.ToLower(getEnv("SEED_DATA", "true")) == "true"

	// Parse CORS origins from comma-separated string
	corsOrigins := strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"), ",")
	// Trim whitespace from each origin
	for i, origin := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		// Server
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "research_paper_analyzer"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-in-production"),
		JWTExpiryHours: jwtExpiry,

		// AI
		AIProvider:   getEnv("AI_PROVIDER", "mock"),
		BedrockModel: getEnv("BEDROCK_MODEL_ID", "anthropic.claude-3-sonnet-20240229-v1:0"),
		AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),

		// Storage
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),
		S3BucketName:    getEnv("S3_BUCKET_NAME", "research-paper-analyzer-uploads"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),

		// CORS
		CORSAllowedOrigins: corsOrigins,

		// Seed
		SeedData: seedData,
	}
}

// getEnv retrieves an environment variable by key, returning a fallback value
// if the variable is not set or is empty.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
