package data

import (
	//"context"
	//"database/sql"
	//	"github.com/cohune-cabbage/di/internal/validator"
	"time"
)

type Login struct {
	ID int64 `json:"id"`
	//	email     string    `json:"email"`
	//	password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}
