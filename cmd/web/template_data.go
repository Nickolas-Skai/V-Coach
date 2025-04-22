package main

import (
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
	CurrentQuestion       string
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
	Question              *data.QuestionData
	QuestionJSON          string
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}

func NewHomePageData() *HomePageData {
	return &HomePageData{
		Title:       "Default Title",
		Header:      "Default Header",
		Description: "Default Description",
	}
}
