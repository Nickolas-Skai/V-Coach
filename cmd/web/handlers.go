package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	//"github.com/cohune-cabbage/di/internal/validator"
	"strconv"

	"github.com/gorilla/csrf"

	"encoding/gob"
	"github.com/cohune-cabbage/di/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gob.Register(&data.Login{})   // Register individual struct
	gob.Register([]*data.Login{}) // Register slices if storing multiple logins
	gob.Register(&data.User{})    // Optional: any other types stored in session
}
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
	//load the user id and role from the session
	userID := app.sessionManager.GetInt(r, "user_id")
	role := app.sessionManager.GetString(r, "user_role")
	app.logger.Info("User ID", "user_id", userID)
	app.logger.Info("User Role", "user_role", role)

}

func (app *application) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in
	if app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Prepare template data
	data := app.addDefaultData(&TemplateData{
		Title:           "Login",
		HeaderText:      "Login to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
		NavLogo:         "static/images/logo.svg",
		CSRFToken:       template.JS(csrf.TemplateField(r)),
	}, w, r)

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

	// Authenticate the user using the new query
	user, err := app.loginModel.Authenticate(email, password)
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

	// Store user ID and role in the session
	err = app.sessionManager.Put(r, w, "user_id", user)
	if err != nil {
		app.logger.Error("Failed to store user ID in session", "error", err)
		app.serverError(w, err)
		return
	}
	app.logger.Info("User ID stored in session", "user_id", user)
	// Set authenticated to true
	err = app.sessionManager.Put(r, w, "IsAuthenticated", true)
	if err != nil {
		app.logger.Error("Failed to store IsAuthenticated in session", "error", err)
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
	//getuserrolebyid
	role, err := app.loginModel.GetUserRoleByID(int64(user))
	if err != nil {
		app.logger.Error("Failed to get user role", "error", err)
		app.serverError(w, err)
		return
	}
	// Store user role in the session
	err = app.sessionManager.Put(r, w, "user_role", role)
	if err != nil {
		app.logger.Error("Failed to store user role in session", "error", err)
		app.serverError(w, err)
		return
	}
	app.logger.Info("User role stored in session", "user_role", role)

	// Redirect to the appropriate dashboard based on the user role

	http.Redirect(w, r, "/", http.StatusSeeOther)

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
		app.logger.Error("Failed to fetch schools from database", "error", err)
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
		err := rows.Scan(&school.ID, sql.RawBytes{}, &school.Name)
		if err != nil {
			app.logger.Error("Failed to scan school row", "error", err)
			app.serverError(w, err)
			return
		}
		schools = append(schools, school)
	}

	if len(schools) == 0 {
		app.logger.Warn("No schools found in the database")
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
func (app *application) CoachDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := app.addDefaultData(&TemplateData{
		Title:           "Coach Dashboard",
		HeaderText:      "Welcome to Your Dashboard",
		PageDescription: "Manage your coaching sessions here",
		NavLogo:         "static/images/logo.svg",
	}, w, r)

	err := app.render(w, http.StatusOK, "coach_dashboard.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render template", "template", "coach_dashboard.tmpl", "error", err)
		app.serverError(w, err)
	}
}

// handlers.go (InterviewHandler with embedded questions and logic)
func (app *application) InterviewHandler(w http.ResponseWriter, r *http.Request) {

	//call createinterviewsession function if not already created and link it to the user
	// Check if the user is logged in
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//is authenticated to get userid

	if !app.sessionManager.Exists(r, "IsAuthenticated") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//log what isAuthenticated is having
	app.logger.Info("IsAuthenticated", "IsAuthenticated", app.sessionManager.Exists(r, "IsAuthenticated"))

	//says its zero
	app.logger.Info("User ID", "user_id", app.sessionManager.GetInt(r, "user_id"))
	// Check if the session ID exists in the session
	TeacherID := app.sessionManager.GetInt(r, "user_id")
	app.logger.Info("Teacher ID", "user_id", TeacherID)

	sessionID := app.sessionManager.GetInt(r, "session_id")
	if sessionID == 0 {
		// Create a new interview session
		sessionID, err := app.InterviewResponseModel.CreateInterviewSession(TeacherID)
		if err != nil {
			app.logger.Error("Failed to create interview session", "error", err)
			app.serverError(w, err)
			return
		}
		app.sessionManager.Put(r, w, "session_id", sessionID)
		app.logger.Info("Created new interview session", "session_id", sessionID)
		app.logger.Info("User ID", "user_id", app.sessionManager.GetInt(r, "user_id"))
		app.logger.Info("Session ID", "session_id", sessionID)
		app.logger.Info("Teacher ID", "user_id", TeacherID)
	}

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

	// Add default data to the template
	data := app.addDefaultData(&TemplateData{
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
	}, w, r)

	// Store state in session
	app.sessionManager.Put(r, w, "current_question_index", qIndex)
	app.sessionManager.Put(r, w, "current_question_id", current.ID)
	app.sessionManager.Put(r, w, "current_question_type", current.Type)
	app.sessionManager.Put(r, w, "current_question_options", current.Options)
	app.sessionManager.Put(r, w, "current_question_required", current.Required)
	app.sessionManager.Put(r, w, "current_question_text", current.Text)

	// Prepare template data
	data = &TemplateData{
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

// interviewsuccess
func (app *application) InterviewSuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Prepare template data
	data := app.addDefaultData(&TemplateData{
		Title:           "Interview Completed",
		HeaderText:      "Thank You for Completing the Interview",
		PageDescription: "Your responses have been recorded.",
		NavLogo:         "static/images/logo.svg",
	}, w, r)

	err := app.render(w, http.StatusOK, "Interview_success.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render template", "template", "Interview_success.tmpl", "url", r.URL, "method", r.Method, "error", err)
		app.serverError(w, err)
	}
}

// Next question handler
func (app *application) NextQuestionHandler(w http.ResponseWriter, r *http.Request) {

	app.logger.Info("Received request for NextQuestionHandler", "method", r.Method, "url", r.URL.String())

	sessionID := app.sessionManager.GetInt(r, "session_id")
	if sessionID == 0 {
		app.logger.Warn("Missing session_id in session")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // Increase file upload limit
		app.logger.Error("Failed to parse multipart form", "error", err)
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}

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

	questions, err := app.questionModel.GetActiveQuestions()
	if err != nil {
		app.logger.Error("Failed to fetch questions", "error", err)
		http.Error(w, "Failed to load questions", http.StatusInternalServerError)
		return
	}
	if len(questions) == 0 {
		app.logger.Warn("No questions available")
		http.Error(w, "No questions available", http.StatusNotFound)
		return
	}

	if qIndex >= len(questions) {
		http.Redirect(w, r, "/interview/complete", http.StatusSeeOther)
		return
	}

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
		file, header, err := r.FormFile("answer")
		if err != nil {
			app.logger.Warn("Failed to retrieve file", "error", err)
		} else if file == nil {
			app.logger.Warn("No file uploaded")
		} else {
			defer file.Close()

			// Define the directory to save the file
			uploadDir := "./uploads/"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				app.logger.Error("Failed to create upload directory", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Generate the file path
			filePath := uploadDir + header.Filename

			// Save the file to the server
			outFile, err := os.Create(filePath)
			if err != nil {
				app.logger.Error("Failed to save file", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, file)
			if err != nil {
				app.logger.Error("Failed to write file to disk", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			app.logger.Info("File uploaded and saved", "filename", header.Filename, "path", filePath)
			answer = filePath // Save the file path as the answer
			isAnswered = true
		}
	default:
		isAnswered = answer != ""
	}

	if !isAnswered {
		app.logger.Warn("No answer submitted", "question_id", currentQ.ID)
		http.Error(w, "Please answer before continuing", http.StatusBadRequest)
		return
	}

	app.logger.Info("Saving response", "session_id", sessionID, "question_id", currentQ.ID, "answer", answer)

	err = app.InterviewResponseModel.SaveAnswer(&data.InterviewResponse{
		SessionID:   sessionID,
		QuestionID:  currentQ.ID,
		Answer:      answer,
		SubmittedAt: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		app.logger.Error("Failed to save interview response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	app.logger.Info("Saved answer", "question_id", currentQ.ID, "answer", answer)

	nextIndex := qIndex + 1
	if nextIndex >= len(questions) {
		http.Redirect(w, r, "/interview/complete", http.StatusSeeOther)
		return
	}

	nextQ := questions[nextIndex]
	data := app.addDefaultData(&TemplateData{
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
	}, w, r)

	app.logger.Info("Rendering template with data", "data", data)

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

// handler to get all a list of all the teachers
func (app *application) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch all teachers from the database
	teachers, err := app.UserModel.GetAllTeachers()
	if err != nil {
		app.logger.Error("Failed to fetch teachers from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := app.addDefaultData(&TemplateData{
		Title:           "All Teachers",
		HeaderText:      "List of All Teachers",
		PageDescription: "Manage your coaching sessions here",
		NavLogo:         "static/images/logo.svg",
	}, w, r)

	data.Teachers = make([]struct {
		ID     int
		Name   string
		Email  string
		Age    int
		School string
	}, len(teachers))
	for i, teacher := range teachers {
		data.Teachers[i] = struct {
			ID     int
			Name   string
			Email  string
			Age    int
			School string
		}{
			ID:    teacher.ID,
			Name:  teacher.Name,
			Email: teacher.Email,
			Age:   teacher.Age,
			School: func() string {
				if teacher.School.Valid {
					return strconv.FormatInt(teacher.School.Int64, 10)
				}
				return ""
			}(),
		}
	}

	err = app.render(w, http.StatusOK, "teacher_list.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render template", "template", "all_teachers.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// teacher details

func (app *application) TeacherDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/coach/teachers/")
	idStr = strings.TrimSuffix(idStr, "/details")

	app.logger.Info("trimmed successfully", "idStr", idStr)

	// Validate the teacher ID
	if idStr == "" {
		app.logger.Warn("Invalid teacher ID")
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(idStr)
	if err != nil || teacherID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Fetch teacher details by ID
	teacher, err := app.UserModel.GetTeacherByID(teacherID)
	if err != nil {
		app.logger.Error("Failed to fetch teacher details", "teacherID", teacherID, "error", err)
		app.serverError(w, err)
		return
	}
	app.logger.Info("Fetched teacher details", "teacherID", teacherID)
	//using school id to get the school name
	schoolName, err := app.UserModel.GetSchoolNameByID(int(teacher.School.Int64))
	if err != nil {
		app.logger.Error("Failed to fetch school name", "schoolID", teacher.School.Int64, "error", err)
		app.serverError(w, err)
		return
	}
	if err != nil {
		app.logger.Error("Failed to fetch school name", "schoolID", teacher.School.Int64, "error", err)
		app.serverError(w, err)
		return
	}

	app.logger.Info("Fetched school details", "schoolID", teacher.School.Int64)

	// Prepare template data with the name of the school
	data := &TemplateData{
		Title: "Teacher Details",
		Teachers: []struct {
			ID     int
			Name   string
			Email  string
			Age    int
			School string
		}{
			{
				ID:     teacher.ID,
				Name:   teacher.Name,
				Email:  teacher.Email,
				Age:    teacher.Age,
				School: strconv.Itoa(schoolName),
			},
		},
	}

	// Render the template
	err = app.render(w, http.StatusOK, "teacher_details.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render template", "template", "teacher_details.tmpl", "error", err)
		app.serverError(w, err)
		return
	}
}

//teacher interview list handler

func (app *application) AllInterviewSessionsListbyTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !app.sessionManager.Exists(r, "user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch all teachers from the database
	teacherID := app.sessionManager.GetInt(r, "user_id")
	if teacherID == 0 {
		app.logger.Warn("Missing teacher ID in session")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//get id from the url
	idStr := strings.TrimPrefix(r.URL.Path, "/coach/teachers/")
	idStr = strings.TrimSuffix(idStr, "/sessions")
	app.logger.Info("trimmed successfully", "idStr", idStr)
	// Validate the teacher ID
	if idStr == "" {
		app.logger.Warn("Invalid teacher ID")
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	teacherID, err := strconv.Atoi(idStr)
	if err != nil || teacherID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	sessions, err := app.InterviewResponseModel.GetAllSessionsByTeacherID(teacherID)
	if err != nil {
		app.logger.Error("Failed to fetch teacher from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	app.logger.Info("Fetched interview sessions", "teacherID", teacherID)
	// Check if there are no sessions
	if len(sessions) == 0 {
		app.logger.Warn("No interview sessions found for the teacher")
		http.Error(w, "No interview sessions found", http.StatusNotFound)
		return
	}

	// Prepare template data
	data := app.addDefaultData(&TemplateData{
		Title:           "Ale",
		HeaderText:      fmt.Sprintf("List of all the sessions made by %s", "Teacher Name"),
		PageDescription: "Manage your coaching sessions here",
		NavLogo:         "static/images/logo.svg",
	}, w, r)
	data.Sessions = make([]struct {
		ID        int
		TeacherID int
		StartTime time.Time
	}, len(sessions))
	for i, session := range sessions {
		data.Sessions[i] = struct {
			ID        int
			TeacherID int
			StartTime time.Time
		}{
			ID:        session,    // Assuming session is an int representing the ID
			TeacherID: teacherID,  // Use teacherID from the context
			StartTime: time.Now(), // Replace with actual start time if available
		}
	}
	err = app.render(w, http.StatusOK, "interview_sessions_list_by_teacher.tmpl", data)
	if err != nil {
		app.logger.Error("Failed to render template", "template", "interview_sessions_list.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
