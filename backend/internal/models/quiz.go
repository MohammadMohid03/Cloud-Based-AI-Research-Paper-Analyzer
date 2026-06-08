package models

import (
	"time"
)

// QuizQuestion represents a multiple-choice question generated from a paper.
// Maps to the "quiz_questions" table in PostgreSQL.
type QuizQuestion struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	PaperID       uint      `json:"paper_id" gorm:"index;not null"`          // Which paper this quiz belongs to
	Question      string    `json:"question" gorm:"type:text;not null"`      // The question text
	OptionA       string    `json:"option_a" gorm:"type:text;not null"`      // Option A
	OptionB       string    `json:"option_b" gorm:"type:text;not null"`      // Option B
	OptionC       string    `json:"option_c" gorm:"type:text;not null"`      // Option C
	OptionD       string    `json:"option_d" gorm:"type:text;not null"`      // Option D
	CorrectAnswer string    `json:"correct_answer" gorm:"type:varchar(1);not null"` // "A", "B", "C", or "D"
	Explanation   string    `json:"explanation" gorm:"type:text"`            // Why the answer is correct
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// QuizQuestionAI is used to receive quiz data from the AI service
// before storing it in the database.
type QuizQuestionAI struct {
	Question      string `json:"question"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
}
