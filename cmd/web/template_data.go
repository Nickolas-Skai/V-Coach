package main

import "github.com/cohune-cabbage/di/internal/data"

type TemplateData struct {
	Title           string
	HeaderText      string
	PageDescription string
	NavLogo         string
	FormErrors      map[string]string
	FormData        map[string]string
	SignUp          *data.SignUp
	SignUpErrors    map[string]string
	SignUpFormData  map[string]string
	Login           *data.Login
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
