package main

import (
	"database/sql"
	"net/http"

	"github.com/cohune-cabbage/di/internal/data"
	"github.com/cohune-cabbage/di/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) homepage(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:           "Home",
		HeaderText:      "Welcome to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
	}
	err := app.render(w, http.StatusOK, "homepage.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Sign Up",
		}
		err := app.render(w, http.StatusOK, "signup.tmpl", data)
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

	signUp := &data.SignUp{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
		Role:  r.PostForm.Get("role"),
	}

	v := validator.NewValidator()
	v.Check(validator.NotBlank(signUp.Name), "name", "Name cannot be blank")
	v.Check(validator.NotBlank(signUp.Email), "email", "Email cannot be blank")
	v.Check(validator.IsValidEmail(signUp.Email), "email", "Invalid email address")
	//if user already signed up

	if !v.ValidData() {
		data := &TemplateData{
			Title:        "Sign Up",
			SignUp:       signUp,
			SignUpErrors: v.Errors,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	// Save the signUp data to the database (omitted for brevity)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {
	// Example: get the first active question for the logged-in teacher
	userID, err := app.getUserID(r.Header.Get("User-ID")) // Assuming the user ID is passed in the "User-ID" header
	if err != nil {
		app.serverError(w, err)
		return
	}
	//no need for getting role because we are assuming the user is a teacher beause all users are teachers
	question, err := app.GetNextQuestionForTeacher(r.Context(), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create a JSON-serializable struct for the template
	q := &data.QuestionData{
		ID:       question.ID,
		Text:     question.Text,
		Type:     question.Type,
		Required: question.Required,
	}

	if question.Type == "checkbox" || question.Type == "radio" || question.Type == "scale" {
		q.Options = question.Options
	}
	// Render the template with the question data
	data := &TemplateData{
		Title:    "Interview",
		Question: q,
	}
	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}
func (app *application) verifyUserCredentials(email, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err := app.db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil // Password mismatch
	}

	return true, nil
}

func (app *application) managequestionsHandler(w http.ResponseWriter, r *http.Request) {

	convertToMapInterface := func(input map[string]string) map[string]interface{} {
		output := make(map[string]interface{})
		for key, value := range input {
			output[key] = value
		}
		return output
	}
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Manage Questions",
		}
		err := app.render(w, http.StatusOK, "managequestions.tmpl", data)
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

	qd := &data.QuestionData{
		Text:                  r.PostForm.Get("question"),
		Type:                  r.PostForm.Get("question_type"),
		AllowConfidenceRating: r.PostForm.Get("allow_confidence_rating") == "on",
	}

	v := validator.NewValidator()
	v.Check(validator.NotBlank(qd.Text), "question", "Question cannot be blank")
	v.Check(validator.NotBlank(qd.Type), "question_type", "Question type cannot be blank")
	v.Check(validator.IsValidQuestionType(qd.Type), "question_type", "Invalid question type")

	options :=
		r.PostForm["options"] // Assuming options are sent as a slice of strings
	if len(options) > 0 {	
		qd.Options = make([]string, len(options))
		for i, option := range options {
			qd.Options[i] = option
		}
	}
	v.Check(validator.MinLength(qd.Text, 5), "question", "Question must be at least 5 characters long")
	v.Check(validator.MaxLength(qd.Text, 100), "question", "Question must be at most 100 characters long")
	v.Check(validator.MinLength(qd.Type, 3), "question_type", "Question type must be at least 3 characters long")
	v.Check(validator.MaxLength(qd.Type, 20), "question_type", "Question type must be at most 20 characters long")
	if !v.ValidData() {
		data := &TemplateData{
			Title:        "Manage Questions",
			QuestionData: qd,
			Data:         convertToMapInterface(v.Errors),
		}
		err := app.render(w, http.StatusUnprocessableEntity, "managequestions.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	// Save the question data to the database (omitted for brevity)
	// For example, you might have a function like this:
	err = app.questionModel.Insert(r.Context(), qd)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/interview", http.StatusSeeOther)
}

func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Interview",
	}
	err := app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) CoachDashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Coach Dashboard",
	}
	err := app.render(w, http.StatusOK, "coach_dashboard.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) ManageQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Manage Questions",
	}
	err := app.render(w, http.StatusOK, "manage_questions.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) EditQuestionHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Edit Question",
	}
	err := app.render(w, http.StatusOK, "edit_question.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) NewQuestionHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "New Question",
	}
	err := app.render(w, http.StatusOK, "new_question.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) TeacherSessionsHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Teacher Sessions",
	}
	err := app.render(w, http.StatusOK, "teacher_sessions.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title: "Delete Teacher",
	}
	err := app.render(w, http.StatusOK, "delete_teacher.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}
