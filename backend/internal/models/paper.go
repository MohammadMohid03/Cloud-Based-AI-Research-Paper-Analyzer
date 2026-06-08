package models

import (
	"time"
)

// PaperStatus represents the processing state of a paper.
type PaperStatus string

const (
	StatusPending    PaperStatus = "pending"    // Paper uploaded, awaiting processing
	StatusProcessing PaperStatus = "processing" // AI analysis in progress
	StatusCompleted  PaperStatus = "completed"  // Analysis finished successfully
	StatusFailed     PaperStatus = "failed"     // Analysis encountered an error
)

// Paper represents an uploaded research paper.
// Maps to the "papers" table in PostgreSQL.
type Paper struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"index;not null"`                         // Foreign key to users table
	Title     string      `json:"title" gorm:"type:varchar(500);not null"`               // Paper title
	FileURL   string      `json:"file_url" gorm:"type:varchar(1000)"`                    // URL or local path to the PDF
	S3Key     string      `json:"s3_key,omitempty" gorm:"type:varchar(500)"`             // S3 object key (if stored in S3)
	Status    PaperStatus `json:"status" gorm:"type:varchar(20);default:'pending';not null"` // Processing status
	RawText   string      `json:"raw_text,omitempty" gorm:"type:text"`                   // Extracted text from PDF
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Analysis      *PaperAnalysis  `json:"analysis,omitempty" gorm:"foreignKey:PaperID;constraint:OnDelete:CASCADE"`
	Chunks        []PaperChunk    `json:"chunks,omitempty" gorm:"foreignKey:PaperID;constraint:OnDelete:CASCADE"`
	ChatHistories []ChatHistory   `json:"chat_histories,omitempty" gorm:"foreignKey:PaperID;constraint:OnDelete:CASCADE"`
	QuizQuestions []QuizQuestion  `json:"quiz_questions,omitempty" gorm:"foreignKey:PaperID;constraint:OnDelete:CASCADE"`
}

// PaperResponse is a simplified response for listing papers (without heavy fields).
type PaperResponse struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	Status    PaperStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ToPaperResponse converts a Paper to a lightweight PaperResponse.
func (p *Paper) ToPaperResponse() PaperResponse {
	return PaperResponse{
		ID:        p.ID,
		Title:     p.Title,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
