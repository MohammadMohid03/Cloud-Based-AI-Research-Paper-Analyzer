package models

import (
	"time"
)

// ChatHistory stores a single question-answer exchange between a user and the AI
// about a specific paper. Maps to the "chat_histories" table in PostgreSQL.
type ChatHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PaperID   uint      `json:"paper_id" gorm:"index;not null"`     // Which paper the chat is about
	UserID    uint      `json:"user_id" gorm:"index;not null"`      // Which user asked the question
	Question  string    `json:"question" gorm:"type:text;not null"` // User's question
	Answer    string    `json:"answer" gorm:"type:text;not null"`   // AI's response
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ChatRequest represents the expected JSON body for a chat message.
type ChatRequest struct {
	Question string `json:"question" binding:"required,min=1"`
}

// ChatResponse is returned after the AI answers a question.
type ChatResponse struct {
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatMessage is used internally to pass chat history context to the AI service.
type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // Message content
}
