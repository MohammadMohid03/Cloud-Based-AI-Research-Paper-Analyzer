package models

import (
	"time"
)

// PaperAnalysis holds the AI-generated analysis results for a paper.
// Maps to the "paper_analyses" table in PostgreSQL.
type PaperAnalysis struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PaperID     uint      `json:"paper_id" gorm:"uniqueIndex;not null"`   // One analysis per paper
	Summary     string    `json:"summary" gorm:"type:text"`               // AI-generated summary
	KeyFindings string    `json:"key_findings" gorm:"type:text"`          // Bullet-point key findings
	Methodology string    `json:"methodology" gorm:"type:text"`           // Research methodology description
	Limitations string    `json:"limitations" gorm:"type:text"`           // Identified limitations
	FutureScope string    `json:"future_scope" gorm:"type:text"`          // Future research directions
	Keywords    string    `json:"keywords" gorm:"type:text"`              // Comma-separated keywords
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// AnalysisReport is used for the downloadable JSON report.
type AnalysisReport struct {
	PaperTitle  string   `json:"paper_title"`
	Summary     string   `json:"summary"`
	KeyFindings []string `json:"key_findings"`
	Methodology string   `json:"methodology"`
	Limitations []string `json:"limitations"`
	FutureScope []string `json:"future_scope"`
	Keywords    []string `json:"keywords"`
	GeneratedAt string   `json:"generated_at"`
}
