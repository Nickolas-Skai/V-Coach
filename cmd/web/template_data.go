package main

import (
	"html/template"

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
	Questions             []*data.QuestionData // Add this field for questions
	QuestionsJSON         template.JS          // Add this field for JSON representation of questions
	InterviewResponse     *data.InterviewResponseModel
	InterviewResponseJSON *data.InterviewResponseModel
	IsLoggedIn            bool
	NextURL               string
	ShowNextButton        bool
	PreviousURL           string
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
