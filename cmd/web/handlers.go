package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"github.com/cohune-cabbage/di/internal/validator"
	"strconv"

	"github.com/cohune-cabbage/di/internal/data"
)

func (app *application) serverError(w http.ResponseWriter, _ error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// to get to pages
func (app *application) homepage(w http.ResponseWriter, r *http.Request) {

	data := &TemplateData{
		Title:           "Home",
		HeaderText:      "Welcome to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
		Greeting:        "",
	}
	err := app.render(w, http.StatusOK, "homepage.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "homepage.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}

}

func (app *application) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:           "Login",
		HeaderText:      "Login to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
	}
	err := app.render(w, http.StatusOK, "login.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "login.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}
func (app *application) ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form data", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create a new Login struct and populate it with the form data
	login := &data.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	// Validate the login data
	err = app.loginModel.ValidateLogin(login)
	if err != nil {
		app.logger.Error("failed to validate login data", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the user exists in the database
	user, err := app.loginModel.GetUserByEmail(login.Email)
	if err != nil {
		app.logger.Error("failed to get user by email", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		app.logger.Error("user not found", "url", r.URL, "method", r.Method)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	// Check if the password is correct
	err = app.loginModel.CheckPassword(user, login.Password)
	if err != nil {
		app.logger.Error("failed to check password", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Set a cookie with the user's name
	http.SetCookie(w, &http.Cookie{
		Name:  "user_name",
		Value: user.Email,
	})
	// Redirect to the dashboard after successful login and is a coach if not redirect to the homepage with name greeting
	if user.Role == "coach" {
		http.Redirect(w, r, "/coach_dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func (app *application) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:           "Sign Up",
		HeaderText:      "Create an Account",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
	}
	err := app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "signup.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}
func (app *application) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	parseIntPointer := func(value string) *int {
		if value == "" {
			return nil
		}
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}
		return &parsed
	}
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// Remove misplaced lines and ensure proper struct initialization
	}
	// Create a new SignUp struct and populate it with the form data
	signup := &data.SignUp{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
		Age:      parseIntPointer(r.FormValue("age")),
		SchoolID: parseIntPointer(r.FormValue("school_id")),
		CoachID:  parseIntPointer(r.FormValue("coach_id")),
	}
	// Validate the signup data
	err = app.signUpModel.ValidateSignUp(signup)
	if err != nil {
		app.logger.Error("failed to validate signup data", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert the new user into the database
	err = app.signUpModel.InsertUser(signup)
	if err != nil {
		app.logger.Error("failed to insert user into database", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
		return
	}
	// Redirect to the login page after successful signup
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) CoachDashBoardHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:           "Coach Dashboard",
		HeaderText:      "Welcome to Your Dashboard",
		PageDescription: "Your virtual coaching assistant.",

		NavLogo: "ui/static/images/logo.svg",
	}
	err := app.render(w, http.StatusOK, "coach_dashboard.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "Coach_dashboard.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}

// handlers.go (InterviewHandler with embedded questions and logic)
// InterviewHandler with dynamic question serving
func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.logger.Error("Failed to fetch interview questions", "error", err)
		app.serverError(w, err)
		return
	}

	if len(questions) == 0 {
		app.logger.Warn("No questions available for the interview")
		http.Error(w, "No interview questions available", http.StatusNotFound)
		return
	}

	// Read the current question index from the query param
	qIndexStr := r.URL.Query().Get("q")
	qIndex := 0
	if qIndexStr != "" {
		parsed, err := strconv.Atoi(qIndexStr)
		if err == nil && parsed >= 0 && parsed < len(questions) {
			qIndex = parsed
		}
	}

	// Pick only the current question
	q := questions[qIndex]
	questionData := &data.QuestionData{
		ID:   q.ID,
		Text: q.Text,
		Type: q.Type,
		Options: func() []string {
			optionsJSON, _ := json.Marshal(q.Options)
			var options []string
			json.Unmarshal(optionsJSON, &options)
			return options
		}(),
		Required: q.Required,
	}

	data := &TemplateData{
		Title:     fmt.Sprintf("Interview Question %d", qIndex+1),
		Questions: []*data.QuestionData{questionData},
	}

	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render interview question", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func (app *application) SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {

}
