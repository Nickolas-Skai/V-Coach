package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"github.com/cohune-cabbage/di/internal/validator"
	"strconv"

	"github.com/cohune-cabbage/di/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) serverError(w http.ResponseWriter, _ error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// to get to pages
func (app *application) homepage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := app.sessionManager.Exists(r, "user_id")
	userRole := app.sessionManager.GetString(r, "user_role")

	data := &TemplateData{
		Title:           "Home",
		HeaderText:      "Welcome to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
		Greeting:        "",
		IsLoggedIn:      isLoggedIn,
		UserRole:        userRole,
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
func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	if email == "" || password == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	login := &data.Login{
		Email:    email,
		Password: password,
	}
	// Validate the login credentials
	err = app.loginModel.ValidateLogin(login)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check if the user exists in the database
	userID, err := app.loginModel.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.serverError(w, err)
		return
	}
	// Set the session data
	app.sessionManager.Put(r, w, "user_id", userID)
	// Set the session data
	app.sessionManager.Put(r, w, "user_id", login.ID)
	app.sessionManager.Put(r, w, "user_role", login.Role)
	app.sessionManager.Put(r, w, "flash", "Login successful")
	// Redirect to the dashboard or homepage
	if login.Role == "coach" {
		http.Redirect(w, r, "/coach/dashboard", http.StatusSeeOther)
	} else if login.Role == "teacher" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	// Render the login page
	err = app.render(w, http.StatusOK, "login.tmpl", nil)
	if err != nil {
		app.logger.Error("failed to render template", "template", "login.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}
func (app *application) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the list of schools

	rows, err := app.db.Query("SELECT id, name FROM schools ORDER BY name")
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var schools []struct {
		ID   int
		Name string
	}
	for rows.Next() {
		var school struct {
			ID   int
			Name string
		}
		err := rows.Scan(&school.ID, &school.Name)
		if err != nil {
			app.serverError(w, err)
			return
		}
		schools = append(schools, school)
	}

	// Pass the schools to the template
	data := &TemplateData{
		Title:           "Sign Up",
		HeaderText:      "Create an Account",
		PageDescription: "Join V-Coach today.",
		NavLogo:         "static/images/logo.svg",
		Schools:         schools, // Add the list of schools
	}
	err = app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "signup.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}
func (app *application) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	role := r.Form.Get("role")
	ageStr := r.Form.Get("age")
	schoolIDStr := r.Form.Get("school_id")

	// Validate required fields
	if name == "" || email == "" || password == "" || role == "" || ageStr == "" || schoolIDStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//hash the password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Convert age and school_id to integers
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		fmt.Println("Error converting age:", err)
		return
	}

	schoolID, err := strconv.Atoi(schoolIDStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		fmt.Println("Error converting schoolID:", err)
		return
	}

	// Handle coachID as optional
	var coachID *int
	coachIDStr := r.Form.Get("coach_id")
	if coachIDStr != "" {
		coachIDVal, err := strconv.Atoi(coachIDStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			fmt.Println("Error converting coachID:", err)
			return
		}
		coachID = &coachIDVal
	}

	signup := &data.SignUp{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Age:      &age,
		SchoolID: &schoolID,
		CoachID:  coachID, // This will be nil if not provided
	}

	err = app.signUpModel.Insert(signup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r, w, "flash", "Sign up successful")

	//if user is entered to db successfully, redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	err = app.render(w, http.StatusOK, "login.tmpl", nil)
	if err != nil {
		app.logger.Error("failed to render template", "template", "login.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}

}
func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Clear the session
	session, err := app.sessionManager.Store.Get(r, "session")
	if err != nil {
		app.serverError(w, err)
		return
	}
	session.Options.MaxAge = -1 // Expire the session
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) CoachDashBoardHandler(w http.ResponseWriter, r *http.Request) {
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := &TemplateData{
		Title:           "Coach Dashboard",
		HeaderText:      "Welcome to Your Dashboard",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "ui/static/images/logo.svg",
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
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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
	fmt.Printf("Rendering template with data: %+v\n", data)

	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render interview question", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func (app *application) SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

}
