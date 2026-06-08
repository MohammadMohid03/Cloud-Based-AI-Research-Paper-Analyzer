// Package models defines the database models (entities) for the application.
// Each model maps to a PostgreSQL table via GORM ORM.
package models

import (
	"time"
)

// User represents a registered user of the application.
// Maps to the "users" table in PostgreSQL.
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"` // "-" hides from JSON output
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships - a user can have many papers and chat histories
	Papers        []Paper       `json:"papers,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ChatHistories []ChatHistory `json:"chat_histories,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// RegisterRequest represents the expected JSON body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// LoginRequest represents the expected JSON body for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after successful login or registration.
type AuthResponse struct {
	Token string `json:"token"`          // JWT token
	User  User   `json:"user"`           // User details (password hash excluded via json:"-")
}
