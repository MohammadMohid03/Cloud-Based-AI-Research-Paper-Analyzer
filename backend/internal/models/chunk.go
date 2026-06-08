package models

import (
	"time"
)

// PaperChunk represents a segment of extracted text from a paper.
// Papers are split into chunks for efficient context retrieval during chat.
// Maps to the "paper_chunks" table in PostgreSQL.
type PaperChunk struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PaperID    uint      `json:"paper_id" gorm:"index;not null"`       // Foreign key to papers table
	ChunkText  string    `json:"chunk_text" gorm:"type:text;not null"` // The text content of this chunk
	ChunkIndex int       `json:"chunk_index" gorm:"not null"`          // Order index (0-based)
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
