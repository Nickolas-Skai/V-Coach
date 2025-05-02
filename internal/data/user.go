package data

import (
	"database/sql"

	"github.com/cohune-cabbage/di/internal/validator"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	Age      int
	School   sql.NullInt64
	Role     string // "teacher" or "coach"
}

type UserModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}
