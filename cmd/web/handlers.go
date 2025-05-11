package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"github.com/cohune-cabbage/di/internal/validator"
	"strconv"

	"github.com/gorilla/csrf"

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

	data := app.addDefaultData(&TemplateData{
		Title:           "V-Coach",
		HeaderText:      "Welcome to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
	}, w, r)
	app.logger.Info("IsAuthenicated", "IsAuthenticated", app.sessionManager.Exists(r, "IsAuthenticated"))
	err := app.render(w, http.StatusOK, "homepage.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "homepage.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}

}

func (app *application) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in
	if app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Prepare template data
	data := &TemplateData{
		Title:           "Login",
		HeaderText:      "Login to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
		CSRFToken:       template.JS(csrf.Token(r)),
	}

	// Render the login page
	err := app.render(w, http.StatusOK, "login.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render login page", "template", "login.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("Failed to parse login form", "error", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.Form.Get("email"))
	password := r.Form.Get("password")

	// Validate the email and password
	if email == "" || password == "" {
		app.logger.Warn("Missing email or password")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Authenticate the user
	userID, err := app.loginModel.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.logger.Warn("Invalid login attempt", "email", email)
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.logger.Error("Error during authentication", "error", err)
		app.serverError(w, err)
		return
	}

	// Fetch the user details
	user, err := app.loginModel.GetUserByID(userID.ID)
	if err != nil {
		app.logger.Error("Failed to fetch user details", "error", err)
		app.serverError(w, err)
		return
	}
	if user == nil {
		app.logger.Warn("User not found", "user_id", userID.ID)
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	// Store user ID and role in the session

	err = app.sessionManager.Put(r, w, "user_id", user.ID)
	if err != nil {
		app.logger.Error("Failed to store user ID in session", "error", err)
		app.serverError(w, err)
		return
	}
	err = app.sessionManager.Put(r, w, "user_role", user.Role)
	if err != nil {
		app.logger.Error("Failed to store user role in session", "error", err)
		app.serverError(w, err)
		return
	}

	// Set the session expiration time
	session, err := app.sessionManager.Store.Get(r, "session")
	if err != nil {
		app.logger.Error("Failed to get session", "error", err)
		app.serverError(w, err)
		return
	}
	session.Options.MaxAge = 3600 // Set session expiration to 1 hour
	err = session.Save(r, w)
	if err != nil {
		app.logger.Error("Failed to save session", "error", err)
		app.serverError(w, err)
		return
	}
	// Redirect to the appropriate dashboard based on the user role
	if user.Role == "coach" {
		http.Redirect(w, r, "/coach_dashboard", http.StatusSeeOther)
	} else if user.Role == "student" {
		http.Redirect(w, r, "/homepage", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.logger.Error("Failed to fetch interview questions", "error", err)
		app.serverError(w, err)
		return
	}
	//log add the question IDs and the texts
	// to the logger
	// for debugging purposes
	for _, question := range questions {
		app.logger.Info("Question details", "ID", question.ID, "Text", question.Text)
	}
	if len(questions) == 0 {
		app.logger.Warn("No questions available for the interview")
		http.Error(w, "No interview questions available", http.StatusNotFound)
		return
	}

	// Initialize "interview_responses" in the session if not already set
	responses := app.sessionManager.Get(r, "interview_responses")
	if err != nil || responses == nil {
		app.sessionManager.Put(r, w, "interview_responses", []*data.InterviewResponse{})
	}

	// Determine current question index
	qIndex := 0
	if qIndexStr := r.URL.Query().Get("q"); qIndexStr != "" {
		if parsed, err := strconv.Atoi(qIndexStr); err == nil && parsed >= 0 && parsed < len(questions) {
			qIndex = parsed
		}
	}

	current := questions[qIndex]

	questionData := &data.QuestionData{
		ID:   current.ID,
		Text: current.Text,
		Type: current.Type,
		Options: func() []string {
			var opts []string
			if current.Options != nil {
				optionsJSON, _ := json.Marshal(current.Options)
				json.Unmarshal(optionsJSON, &opts)
			}
			return opts
		}(),
		Required: current.Required,
	}

	// Store state in session
	app.sessionManager.Put(r, w, "current_question_index", qIndex)
	app.sessionManager.Put(r, w, "current_question_id", current.ID)
	app.sessionManager.Put(r, w, "current_question_type", current.Type)
	app.sessionManager.Put(r, w, "current_question_options", current.Options)
	app.sessionManager.Put(r, w, "current_question_required", current.Required)
	app.sessionManager.Put(r, w, "current_question_text", current.Text)

	// Prepare template data
	data := &TemplateData{
		Title: fmt.Sprintf("Interview Question %d", qIndex+1),
		CurrentQuestion: &Question{
			ID:       questionData.ID,
			Text:     questionData.Text,
			Type:     questionData.Type,
			Options:  questionData.Options,
			Required: questionData.Required,
		},
		TotalQuestions: len(questions),
		CurrentIndex:   qIndex,
	}

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

// function coach dashboard handler
func (app *application) CoachDashboardHandler(w http.ResponseWriter, r *http.Request) {
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

// Next question handler
func (app *application) NextQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging
	app.logger.Info("Received request for NextQuestionHandler", "method", r.Method, "url", r.URL.String())

	// Parse form values
	if err := r.ParseForm(); err != nil {
		app.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}

	// Log form data
	app.logger.Info("Form data received", "form", r.Form)

	// Get current question index from URL query
	qIndex := 0
	if qStr := r.URL.Query().Get("q"); qStr != "" {
		if parsed, err := strconv.Atoi(qStr); err == nil && parsed >= 0 {
			qIndex = parsed
		} else {
			app.logger.Error("Invalid question index", "error", err)
			http.Error(w, "Invalid question index", http.StatusBadRequest)
			return
		}
	}

	// Fetch questions
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil || len(questions) == 0 {
		app.logger.Error("Failed to fetch questions", "error", err)
		http.Error(w, "Failed to load questions", http.StatusInternalServerError)
		return
	}

	// Log the current question index and total questions
	app.logger.Info("Current question index", "qIndex", qIndex, "totalQuestions", len(questions))

	if qIndex >= len(questions) {
		http.Redirect(w, r, "/interview/complete", http.StatusSeeOther)
		return
	}

	// Validate and save current answer
	currentQ := questions[qIndex]
	answer := r.FormValue("answer")
	isAnswered := false

	switch currentQ.Type {
	case "text", "longtext", "radio", "dropdown":
		isAnswered = answer != ""
	case "checkbox":
		selected := r.Form["answer"]
		isAnswered = len(selected) > 0
		answer = strings.Join(selected, ", ")
	case "file":
		file, _, err := r.FormFile("answer")
		isAnswered = err == nil && file != nil
	default:
		isAnswered = answer != ""
	}

	if !isAnswered {
		app.logger.Warn("No answer submitted", "question_id", currentQ.ID)
		http.Error(w, "Please answer before continuing", http.StatusBadRequest)
		return
	}

	err = app.InterviewResponseModel.SaveAnswer(&data.InterviewResponse{

		QuestionID: currentQ.ID,
		Answer:     answer,
	})
	if err != nil {
		app.logger.Error("Failed to save interview response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Move to next question
	nextIndex := qIndex + 1
	if nextIndex >= len(questions) {
		http.Redirect(w, r, "/interview/complete", http.StatusSeeOther)
		return
	}

	nextQ := questions[nextIndex]
	data := &TemplateData{
		Title: fmt.Sprintf("Interview Question %d", nextIndex+1),
		CurrentQuestion: &Question{
			ID:       nextQ.ID,
			Text:     nextQ.Text,
			Type:     nextQ.Type,
			Required: nextQ.Required,
			Options: func() []string {
				var options []string
				b, _ := json.Marshal(nextQ.Options)
				json.Unmarshal(b, &options)
				return options
			}(),
		},
		CurrentIndex:   nextIndex,
		TotalQuestions: len(questions),
	}

	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) PreviousQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current question index from the query parameter
	qIndexStr := r.URL.Query().Get("q")
	qIndex := 0
	if qIndexStr != "" {
		parsed, err := strconv.Atoi(qIndexStr)
		if err == nil && parsed > 0 {
			qIndex = parsed - 1 // Decrement the index for the previous question
		}
	}

	// Fetch all active questions from the database
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.logger.Error("Failed to fetch interview questions", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ensure the index is within bounds
	if qIndex < 0 || qIndex >= len(questions) {
		http.Error(w, "Invalid question index", http.StatusBadRequest)
		return
	}

	// Fetch the previous question
	question := questions[qIndex]
	questionData := &data.QuestionData{
		ID:       question.ID,
		Text:     question.Text,
		Type:     question.Type,
		Options:  question.Options,
		Required: question.Required,
	}

	// Save the current question index in the session
	app.sessionManager.Put(r, w, "current_question_index", qIndex)

	// Render the template with the previous question
	data := &TemplateData{
		Title:        fmt.Sprintf("Interview Question %d", qIndex+1),
		Questions:    []*data.QuestionData{questionData},
		CurrentIndex: qIndex,
	}

	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render interview question", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
