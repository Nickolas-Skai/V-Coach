package main

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/cohune-cabbage/di/internal/data"
)

type TemplateData struct {
	Title                 string
	HeaderText            string
	Greeting              string
	PageDescription       string
	NavLogo               string
	FormErrors            map[string]string
	FormData              map[string]string
	SignUp                *data.SignUp
	SignUpErrors          map[string]string
	SignUpFormData        map[string]string
	Login                 *data.Login
	CurrentQuestion       *Question
	QuestionType          string
	Options               []string
	AllowConfidenceRating bool
	SessionID             string
	TipResponse           string
	QuestionData          *data.QuestionData
	QuestionDataErrors    map[string]string
	QuestionDataFormData  map[string]string
	QuestionDataID        int
	QuestionDataText      string
	QuestionDataType      string
	QuestionDataOptions   []string
	QuestionDataRequired  bool
	QuestionDataModel     *data.QuestionModel
	QuestionDataDB        *data.QuestionModel
	ErrorMessage          string
	Data                  map[string]interface{}
	QuestionsJSON         template.JS // Add this field for JSON representation of questions
	InterviewResponse     *data.InterviewResponseModel
	InterviewResponseJSON *data.InterviewResponseModel
	IsAuthenticated       bool
	NextURL               string
	ShowNextButton        bool
	PreviousURL           string
	SuccessMessage        string
	TotalQuestions        int
	CurrentIndex          int
	Schools               []struct {
		ID   int
		Name string
	}
	UserRole              string
	CurrentUserID         int
	CurrentUserRole       string
	CSRFToken             template.JS
	ErrInvalidCredentials error
	ErrInvalidEmail       error
	Teachers              []struct {
		ID     int
		Name   string
		Email  string
		Age    int
		School string
	}
	CSRFField template.JS
	Sessions  []struct {
		ID        int
		TeacherID int
		StartTime time.Time
	}
	SessionDetails struct {
		ID        int
		TeacherID int
		StartTime time.Time
		EndTime   time.Time
		Duration  time.Duration
		Questions []struct {
			ID      int
			Text    string
			Title   string
			Type    string
			Options []string
		}
	}
	SessionNumber    int
	ParticipantName  string
	ParticipantID    int
	Questions        interface{} // Can hold either []*data.QuestionData or []data.QuestionResponse
	ForcoachSessions []struct {
		ID        int
		TeacherID int
		Title     string
		StartTime time.Time
	}
	Errors     []string          // Added field to store inline error messages
	FormValues map[string]string // Holds form input values for repopulation
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:         "Default Title",
		HeaderText:    "Default HeaderText",
		FormErrors:    map[string]string{},
		FormData:      map[string]string{},
		QuestionsJSON: template.JS("[]"), // Initialize with an empty JSON array as a valid template.JS value
	}
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewHomePageData() *HomePageData {
	return &HomePageData{
		Title:       "Default Title",
		Header:      "Default Header",
		Description: "Default Description",
	}
}

type Question struct {
	ID                    int
	Text                  string
	Type                  string
	Options               []string // Only for checkbox/radio/scale
	AllowConfidenceRating bool
	Required              bool
}

func (app *application) addDefaultData(td *TemplateData, _ http.ResponseWriter, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	td.IsAuthenticated = app.IsAuthenticated(r)
	td.CSRFToken = template.JS(app.sessionManager.GetString(r, "csrf_token"))
	td.UserRole = app.sessionManager.GetString(r, "user_role")
	td.CurrentUserID = app.sessionManager.GetInt(r, "user_id")
	td.CurrentUserRole = app.sessionManager.GetString(r, "user_role")
	td.ErrInvalidCredentials = ErrInvalidCredentials
	td.ErrInvalidEmail = errors.New("invalid email address")
	td.Schools = []struct {
		ID   int
		Name string
	}{
		{ID: 1, Name: "Anglican Primary School"},
		{ID: 2, Name: "Saint Joseph's RC School"},
		{ID: 3, Name: "Belize High School"},
		{ID: 4, Name: "Harmony Government School"},
		{ID: 5, Name: "Sunrise Academy"},
	}
	td.Data = make(map[string]interface{})
	td.Data["user_id"] = td.CurrentUserID
	td.Data["user_role"] = td.CurrentUserRole
	td.Data["csrf_token"] = td.CSRFToken
	td.Data["is_authenticated"] = td.IsAuthenticated
	td.Data["schools"] = td.Schools
	td.Data["questions"] = td.Questions
	td.Data["questions_json"] = td.QuestionsJSON
	td.Data["interview_response"] = td.InterviewResponse
	td.Data["interview_response_json"] = td.InterviewResponseJSON
	td.Data["success_message"] = td.SuccessMessage
	td.Data["error_message"] = td.ErrorMessage
	td.Data["next_url"] = td.NextURL
	td.Data["previous_url"] = td.PreviousURL
	td.Data["show_next_button"] = td.ShowNextButton
	td.Data["current_index"] = td.CurrentIndex
	td.Data["total_questions"] = td.TotalQuestions
	td.Data["current_question"] = td.CurrentQuestion
	td.Data["question_data"] = td.QuestionData
	td.Data["question_data_errors"] = td.QuestionDataErrors
	td.Data["question_data_form_data"] = td.QuestionDataFormData
	td.Data["question_data_id"] = td.QuestionDataID
	td.Data["question_data_text"] = td.QuestionDataText
	td.Data["question_data_type"] = td.QuestionDataType
	td.Data["question_data_options"] = td.QuestionDataOptions
	td.Data["question_data_required"] = td.QuestionDataRequired
	td.Data["question_data_model"] = td.QuestionDataModel

	td.Data["question_data_db"] = td.QuestionDataDB
	td.Data["question_type"] = td.QuestionType

	return td
}
