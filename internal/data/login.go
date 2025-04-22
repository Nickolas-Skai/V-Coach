package data

import (
	"database/sql"
	"github.com/cohune-cabbage/di/internal/validator"
	"time"
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

//validate login
