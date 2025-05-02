package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
	"github.com/lib/pq"
)

// Signup represents a signup request.

type SignUp struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password_hash"`
	Role      string    `json:"role"`
	Age       *int      `json:"age,omitempty"`
	SchoolID  *int      `json:"school_id,omitempty"`
	CoachID   *int      `json:"coach_id,omitempty"`
}

type SignUpModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}

func NewSignUpModel(db *sql.DB) *SignUpModel {
	return &SignUpModel{
		DB:        db,
		Validator: validator.NewValidator(),
	}
}

// Insert inserts a new signup record into the database.
func (m *SignUpModel) Insert(signup *SignUp) error {
	// Insert the signup into the database
	query := `INSERT INTO users (name, email, password_hash, role, age, school_id, coach_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	args := []interface{}{
		signup.Name,
		signup.Email,
		signup.Password,
		signup.Role,
		signup.Age,
		signup.SchoolID,
		sql.NullInt64{}, // Default to NULL for coach_id
		time.Now(),
	}

	// Only include coach_id if it is provided
	if signup.CoachID != nil {
		args[6] = *signup.CoachID
	}

	err := m.DB.QueryRow(query, args...).Scan(&signup.ID)
	if err != nil {
		// Check for specific SQL errors
		if err == sql.ErrNoRows {
			return fmt.Errorf("no rows were returned: %w", err)
		}

		// Check for unique constraint violation (example for PostgreSQL)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation
				return fmt.Errorf("a user with this email already exists: %w", err)
			}
		}

		// General database error
		return fmt.Errorf("failed to insert signup record: %w", err)
	}

	return nil
}

// ValidateSignUp validates the signup data.
func (m *SignUpModel) ValidateSignUp(signup *SignUp) error {
	if m.Validator == nil {
		return fmt.Errorf("validator is not initialized")
	}

	m.Validator.Check(validator.NotBlank(signup.Name), "name", "Name cannot be blank")
	m.Validator.Check(validator.IsValidEmail(signup.Email), "email", "Invalid email address")
	m.Validator.Check(validator.MinLength(signup.Password, 8), "password", "Password must be at least 8 characters long")

	if len(m.Validator.Errors) > 0 {
		return fmt.Errorf("validation errors: %v", m.Validator.Errors)
	}

	return nil
}
