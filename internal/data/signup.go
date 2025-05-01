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
			INSERT INTO users (name, email, password_hash, role, age, school_id, coach_id, created_at)
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

// insert a user into the database
func (m *SignUpModel) InsertUser(signup *SignUp) error {
	// Insert the new user into the database
	query := `
			INSERT INTO users (name, email, password_hash, role, age, school_id, coach_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

	// Execute the query and scan the returned ID into the signup struct
	err := m.DB.QueryRow(
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
