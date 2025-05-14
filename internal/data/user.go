package data

import (
	"database/sql"

	"github.com/Nickolas-Skai/internal/validator"
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

// get all teachers
func (m *UserModel) GetAllTeachers() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, name, email FROM users WHERE role = 'teacher'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []User
	for rows.Next() {
		var teacher User
		if err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Email); err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

// get details of a specific teacher
func (m *UserModel) GetTeacherByID(id int) (*User, error) {
	row := m.DB.QueryRow("SELECT id, name, email, age, school_id FROM users WHERE id = $1 AND role = 'teacher'", id)

	var teacher User
	if err := row.Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Age, &teacher.School); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No teacher found with the given ID
		}
		return nil, err // Some other error occurred
	}

	return &teacher, nil

}

// get school by id
func (m *UserModel) GetSchoolNameByID(id int) (int, error) {
	//select name from schools where id = $1
	row := m.DB.QueryRow("SELECT name FROM schools WHERE id = $1", id)
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No school found with the given ID
		}
		return 0, err // Some other error occurred
	}
	return id, nil

}

//validation for use registration

func ValidateUser( v *validator.Validator,user *UserModel) map[string](string) {

errors:= v.Errors
	v.Check(validator.NotBlank(user.Name), "name", "must be provided")
	v.Check(validator.NotBlank(user.Email), "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(validator.NotBlank(user.Password), "password", "must be provided")
	v.Check(validator.MinLength(user.Password, 8), "password", "must be at least 8 characters long")
	v.Check(validator.MaxLength(user.Password, 100), "password", "must not exceed 100 characters")
	v.Check(validator.NotBlank(user.Role), "role", "must be provided")
	v.Check(validator.Matches(user.Role, validator.RoleRX), "role", "must be a valid role")
	v.Check(validator.NotBlank(user.Age), "age", "must be provided")
	v.Check(validator.MinLength(user.Age, 1), "age", "must be at least 1 characters long")
	v.Check(validator.MaxLength(user.Age, 3), "age", "must not exceed 3 characters")
	v.Check(validator.NotBlank(user.School), "school", "must be provided")
	
	return errors
	
	
}
