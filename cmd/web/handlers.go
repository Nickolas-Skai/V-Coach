package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"github.com/cohune-cabbage/di/internal/validator"
	"strconv"

	"github.com/cohune-cabbage/di/internal/data"
)

func (app *application) serverError(w http.ResponseWriter, _ error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// to get to pages
func (app *application) homepage(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user's name from a cookie (or session)
	cookie, err := r.Cookie("user_name")
	var greeting string
	if err == nil {
		// If the cookie exists, use the user's name for the greeting
		greeting = fmt.Sprintf("Welcome back, %s!", cookie.Value)
	} else {
		// If no cookie is found, use a default greeting
		greeting = "Welcome to V-Coach!"
	}

	// Pass the greeting to the template
	data := &TemplateData{
		Title:           "Home",
		HeaderText:      "Welcome to V-Coach",
		Greeting:        greeting,
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
	}
	err = app.render(w, http.StatusOK, "homepage.tmpl", data)
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
	// Step 1: Load all active questions from the database
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.logger.Error("Failed to fetch interview questions", "error", err)
		app.serverError(w, err)
		return
	}

	// Step 2: Return JSON for dynamic use or render template based on Accept header
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(questions)
		return
	}

	// Step 3: Marshal questions to embed into template if rendering as HTML
	qJSON, err := json.Marshal(questions)
	if err != nil {
		app.logger.Error("Failed to marshal questions to JSON", "error", err)
		app.serverError(w, err)
		return
	}

	// Add console.log statement to log the questions variable
	fmt.Println("Loaded Questions:", string(qJSON))

	// Add a check to handle the case when the questions list is empty
	if len(questions) == 0 {
		app.logger.Warn("No questions available for the interview")
	}

	data := &TemplateData{
		Title:         "Interview",
		QuestionsJSON: template.JS(qJSON),
	}

	err = app.render(w, http.StatusOK, "interview.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render interview page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) SubmitInterviewResponseHandler(w http.ResponseWriter, r *http.Request) {
	var submission struct {
		Answers []struct {
			QuestionID int    `json:"question_id"`
			Answer     string `json:"answer"`
		} `json:"answers"`
	}

	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		app.logger.Error("invalid interview response json", "error", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	for _, ans := range submission.Answers {
		if strings.TrimSpace(ans.Answer) == "" || ans.QuestionID <= 0 {
			app.logger.Error("validation failed", "answer", ans)
			http.Error(w, "Validation error: missing or invalid answer.", http.StatusBadRequest)
			return
		}

		response := &data.InterviewResponse{
			QuestionID: ans.QuestionID,
			Answer:     ans.Answer,
		}
		if err := app.InterviewResponseModel.InsertInterviewResponse(response); err != nil {
			app.logger.Error("failed to insert response", "question_id", ans.QuestionID, "error", err)
			app.serverError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
