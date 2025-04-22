package main

import (
	"net/http"

	"fmt"

	"context"
	"database/sql"
	"github.com/cohune-cabbage/di/internal/data"
	"github.com/cohune-cabbage/di/internal/validator"
	_ "github.com/lib/pq"
)

//getting usersthings
//logining in

func (app *application) getUserID(email string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE email = $1`
	err := app.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// getting user role
func (app *application) getUserRole(email string) (string, error) {
	var role string
	query := `SELECT role FROM users WHERE email = $1`
	err := app.db.QueryRow(query, email).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

// getting user name
func (app *application) getName(email string) (string, error) {
	var name string
	query := `SELECT name FROM users WHERE email = $1`
	err := app.db.QueryRow(query, email).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Login",
		}
		err := app.render(w, http.StatusOK, "login.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	login := &data.Login{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	v := validator.NewValidator()
	v.Check(validator.NotBlank(login.Email), "email", "Email cannot be blank")
	v.Check(validator.IsValidEmail(login.Email), "email", "Invalid email address")
	v.Check(validator.NotBlank(login.Password), "password", "Password cannot be blank")

	if !v.ValidData() {
		data := &TemplateData{
			Title:        "Login",
			ErrorMessage: fmt.Sprintf("%v", v.Errors),
			Data: map[string]interface{}{
				"email":    login.Email,
				"password": login.Password,
			},
		}
		err := app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	id, err := app.getUserID(login.Email)
	if err != nil {
		app.serverError(w, err)
		return
	}

	role, err := app.getUserRole(login.Email)
	if err != nil {
		app.serverError(w, err)
		return
	}

	name, err := app.getName(login.Email)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: fmt.Sprintf("%d", id),
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "user_role",
		Value: role,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "user_name",
		Value: name,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) GetNextQuestionForTeacher(ctx context.Context, teacherID int) (*data.QuestionData, error) {
	query := `
        SELECT id, question_text, question_type, options, allow_confidence_rating
        FROM questions
        WHERE teacher_id = $1 AND is_active = true
        ORDER BY created_at ASC
        LIMIT 1
    `
	row := app.db.QueryRowContext(ctx, query, teacherID)

	var question data.QuestionData
	err := row.Scan(&question.ID, &question.Text, &question.Type, &question.Options, &question.AllowConfidenceRating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active questions found
		}
		return nil, err
	}

	return &question, nil
}

func (app *application) InsertQuestion(ctx context.Context, question *data.QuestionData) (int, error) {
	query := `
		INSERT INTO questions (question_text, question_type, options, allow_confidence_rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := app.db.QueryRowContext(ctx, query, question.Text, question.Type, question.Options, question.AllowConfidenceRating).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *QuestionModel) Insert(ctx context.Context, qd QuestionData) (int, error) {
    query := `
        INSERT INTO questions (question_text, question_type, options, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    var id int
    err := m.DB.QueryRowContext(ctx, query, qd.Text, qd.Type, qd.Options, qd.IsActive).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}