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

// get list of teachers
func (m *UserModel) GetTeachers() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, name, email, password, age, school FROM users WHERE role = 'teacher'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []User
	for rows.Next() {
		var teacher User
		if err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Password, &teacher.Age, &teacher.School); err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}
