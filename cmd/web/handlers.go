package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

// Render the interview response page it should load one question at a time
func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {
	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(questions) == 0 {
		app.logger.Error("No questions found")
		http.Error(w, "No questions available GO HANDLER .", http.StatusInternalServerError)
		return
	}

	qJSON, err := json.Marshal(questions)
	if err != nil {
		app.logger.Error("failed to marshal questions to JSON", "error", err)
		http.Error(w, "Failed to process questions", http.StatusInternalServerError)
		return
	}

	data := &TemplateData{
		Title:           "Interview",
		HeaderText:      "Interview Questions",
		PageDescription: "Answer the following questions.",
		NavLogo:         "static/images/logo.svg",
		QuestionsJSON:   template.JS(qJSON), // Pass the JSON data to the template
	}

	//	Render the interview page with the questions
	err = app.render(w, http.StatusOK, "interview.tmpl", data)

	if err != nil {
		app.logger.Error("failed to render template", "template", "interview.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
		return
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
		app.logger.Error("failed to decode JSON body", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, ans := range submission.Answers {
		response := &data.InterviewResponse{
			QuestionID: ans.QuestionID,
			Text:       ans.Answer,
		}

		err = app.InterviewResponseModel.InsertInterviewResponse(response)
		if err != nil {
			app.logger.Error("failed to insert interview response", "questionID", ans.QuestionID, "url", r.URL, "method", r.Method, "error", err)
			app.serverError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Responses submitted successfully"))
}

func (app *application) GetNextQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form data", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create a new InterviewResponse struct and populate it with the form data
	response := &data.InterviewResponse{
		Text:                  r.FormValue("text"),
		Questiontype:          r.FormValue("type"),
		Options:               r.Form["options"],
		AllowConfidenceRating: r.FormValue("allow_confidence_rating") == "on",
		Required:              r.FormValue("required") == "on",
	}
	// Validate the interview response data
	err = app.InterviewResponseModel.ValidateInterviewResponse(response)
	if err != nil {
		app.logger.Error("failed to validate interview response data", "url", r.URL, "method", r.Method, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
