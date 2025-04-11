package data

import (
	//	"context"
	"database/sql"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
)

// Message represents a message sent by a user. messages can be either a string or a number at this point for testing purposes.

type Response struct {
	ID           int       `json:"id" db:"id"`
	SessionID    int       `json:"session_id" db:"session_id"`
	QuestionID   int       `json:"question_id" db:"question_id"`
	ResponseText *string   `json:"response_text,omitempty" db:"response_text"`
	AudioURL     *string   `json:"audio_url,omitempty" db:"audio_url"`
	Confidence   *int      `json:"confidence,omitempty" db:"confidence"`
	SubmittedAt  time.Time `json:"submitted_at" db:"submitted_at"`
}
type ResponseModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}
