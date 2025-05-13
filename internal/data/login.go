package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cohune-cabbage/di/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

type LoginModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}

func NewLoginModel(db *sql.DB) *LoginModel {
	return &LoginModel{
		DB:        db,
		Validator: validator.NewValidator(),
	}
}

//validate login

func (m *LoginModel) ValidateLogin(login *Login) error {
	v := m.Validator
	v.Check(validator.NotBlank(login.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(login.Email), "email", "must be a valid email address")
	v.Check(validator.NotBlank(login.Password), "password", "must be provided")
	v.Check(validator.MinLength(login.Password, 8), "password", "must be at least 8 characters long")
	v.Check(validator.MaxLength(login.Password, 100), "password", "must not exceed 100 characters")

	if !v.ValidData() {
		return fmt.Errorf("validation errors: %v", v.Errors)
	}
	return nil
}
func (m *LoginModel) Insert(login *Login) error {
	// Insert the login into the database
	query := `INSERT INTO logins (email, password, created_at, role) VALUES ($1, $2, $3, $4) RETURNING id`
	args := []interface{}{login.Email, login.Password, time.Now(), login.Role}
	err := m.DB.QueryRow(query, args...).Scan(&login.ID)
	if err != nil {
		return err
	}
	return nil
}
func (m *LoginModel) GetByEmail(email string) (*Login, error) {
	// Retrieve the login from the database by email
	query := `SELECT id, email, password, created_at, role FROM logins WHERE email = $1`
	row := m.DB.QueryRow(query, email)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No login found with the given email
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}
func (m *LoginModel) Update(login *Login) error {
	// Update the login in the database
	query := `UPDATE users SET email = $1, password_hash = $2, role = $3 WHERE id = $4`
	args := []interface{}{login.Email, login.Password, login.Role, login.ID}
	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *LoginModel) Delete(id int64) error {
	// Delete the login from the database
	query := `DELETE FROM users WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *LoginModel) GetByID(id int64) (*Login, error) {
	// Retrieve the login from the database by ID
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE id = $1`
	row := m.DB.QueryRow(query, id)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No login found with the given ID
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}
func (m *LoginModel) GetAll() ([]*Login, error) {
	// Retrieve all users from the database
	query := `SELECT id, email, password_hash, created_at, role FROM users`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logins []*Login
	for rows.Next() {
		var login Login
		err := rows.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
		if err != nil {
			return nil, err
		}
		logins = append(logins, &login)
	}
	return logins, nil
}
func (m *LoginModel) GetByRole(role string) ([]*Login, error) {
	// Retrieve logins from the database by role
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE role = $1`
	rows, err := m.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logins []*Login
	for rows.Next() {
		var login Login
		err := rows.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
		if err != nil {
			return nil, err
		}
		logins = append(logins, &login)
	}
	return logins, nil
}
func (m *LoginModel) GetByEmailAndPassword(email, password string) (*Login, error) {
	// Retrieve the login from the database by email and password
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE email = $1 AND password = $2`
	row := m.DB.QueryRow(query, email, password)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No login found with the given email and password
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}
func (m *LoginModel) GetByEmailAndRole(email, role string) (*Login, error) {
	// Retrieve the login from the database by email and role
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE email = $1 AND role = $2`
	row := m.DB.QueryRow(query, email, role)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No login found with the given email and role
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}

//get user by email

func (m *LoginModel) GetUserByEmail(email string) (*Login, error) {
	// Retrieve the user from the database by email
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE email = $1`
	row := m.DB.QueryRow(query, email)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with the given email
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}
func (m *LoginModel) CheckPassword(user *Login, password string) error {
	// Check if the provided password matches the stored password
	if user.Password != password {
		return fmt.Errorf("invalid password")
	}
	return nil
}
func (m *LoginModel) GetUserByID(id int64) (*Login, error) {
	// Retrieve the user from the database by ID
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE id = $1`
	row := m.DB.QueryRow(query, id)

	var login Login
	err := row.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with the given ID
		}
		return nil, err // Some other error occurred
	}
	return &login, nil
}

// authenticate user
func (m *LoginModel) Authenticate(email, password string) (int, error) {

	var id int64
	var hashedPassword []byte
	var createdAt time.Time

	query := `SELECT id, password_hash, created_at FROM users WHERE email = $1`
	err := m.DB.QueryRow(query, email).Scan(&id, &hashedPassword, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No user found with the given email
		}
		return 0, err // Some other error occurred
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, fmt.Errorf("invalid password")
	}

	// Return the authenticated user
	return int(id), nil
}

// /get user by role
func (m *LoginModel) GetUserRole(role string) ([]*Login, error) {
	// Retrieve the user's role from the database
	query := `SELECT id, email, password_hash, created_at, role FROM users WHERE role = $1`
	rows, err := m.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logins []*Login
	for rows.Next() {
		var login Login
		err := rows.Scan(&login.ID, &login.Email, &login.Password, &login.CreatedAt, &login.Role)
		if err != nil {
			return nil, err
		}
		logins = append(logins, &login)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logins, nil
}
