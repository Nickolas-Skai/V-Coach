package data

import (
	"database/sql"
	"fmt"
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

// Insert inserts a new signup record into the database.
func (m *SignUpModel) Insert(signup *SignUp) error {
	// Check if the user already exists
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	err := m.DB.QueryRow(checkQuery, signup.Email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("a user with this email already exists")
	}

	// Insert the new user into the database
	query := `
			INSERT INTO users (name, email, password, role, age, school_id, coach_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

	// Execute the query and scan the returned ID into the signup struct
	err = m.DB.QueryRow(
		query,
		signup.Name,
		signup.Email,
		signup.Password,
		signup.Role,
		signup.Age,
		signup.SchoolID,
		signup.CoachID,
		time.Now(),
	).Scan(&signup.ID)

	if err != nil {
		return err
	}

	return nil
}
