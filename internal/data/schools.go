package data

import (
	"database/sql"

	"github.com/cohune-cabbage/di/internal/validator"
)

type SchoolModel struct {
	DB        *sql.DB
	Validator *validator.Validator
}

func NewSchoolModel(db *sql.DB) *SchoolModel {
	return &SchoolModel{
		DB:        db,
		Validator: validator.NewValidator(),
	}
}

// GetSchoolNameByID
func (m *SchoolModel) GetSchools(id int) (string, error) {
	query := `SELECT name FROM schools WHERE id = $1`
	row := m.DB.QueryRow(query, id)
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No school found with the given ID
		}
		return "", err // Other error occurred
	}
	return name, nil
}
