package data

import (
	"database/sql"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
)

// Signup represents a signup request.

type SignUp struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Age       *int      `json:"age,omitempty"`
	SchoolID  *int      `json:"school_id,omitempty"`
	CoachID   *int      `json:"coach_id,omitempty"`
}

type SignUpModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}
