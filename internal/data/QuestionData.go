package data

import (
	//	"context"
	"database/sql"

	"github.com/cohune-cabbage/di/internal/validator"
)

// Question represents a question in the system.
type QuestionData struct {
	ID                    int
	Text                  string
	Type                  string
	Options               []string // Only for checkbox/radio/scale
	AllowConfidenceRating bool
	Required              bool
}
type QuestionModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}
