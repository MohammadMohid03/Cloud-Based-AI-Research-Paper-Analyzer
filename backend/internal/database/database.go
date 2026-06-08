// Package database handles PostgreSQL connection management using GORM.
// It provides initialization and access to the database connection pool.
package database

import (
	"fmt"
	"log"

	"github.com/research-paper-analyzer/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection instance.
// It is initialized once at startup and shared across the application.
var DB *gorm.DB

// Initialize creates and configures the PostgreSQL connection using GORM.
// It sets up connection pooling and configures logging based on the environment.
func Initialize(cfg *config.Config) error {
	// Build the PostgreSQL connection string (DSN)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Configure GORM logger level based on Gin mode
	logLevel := logger.Info
	if cfg.GinMode == "release" {
		logLevel = logger.Warn // Less verbose logging in production
	}

	// Open the database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying *sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings for production performance
	sqlDB.SetMaxIdleConns(10)    // Keep up to 10 idle connections ready
	sqlDB.SetMaxOpenConns(100)   // Allow up to 100 concurrent connections
	// Note: SetConnMaxLifetime can be set if needed for long-running servers

	log.Println("✅ Database connected successfully")
	return nil
}

// GetDB returns the global database connection instance.
// Use this in handlers and services to access the database.
func GetDB() *gorm.DB {
	return DB
}

// Close gracefully closes the database connection.
// Should be called when the application shuts down.
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
