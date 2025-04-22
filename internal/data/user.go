package data

import (
	"database/sql"

	"github.com/cohune-cabbage/di/internal/validator"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Age      int
	School   string
	Role     string // "teacher" or "coach"
}

type UserModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}
