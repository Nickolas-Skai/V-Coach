package main

import (
	"net/http"
)

func (app *application) NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:          "Default Title",
		HeaderText:     "Default HeaderText",
		FormErrors:     map[string]string{},
		FormData:       map[string]string{},
		SignUpErrors:   map[string]string{},
		SignUpFormData: map[string]string{},

		Data:            make(map[string]interface{}),
		ShowNextButton:  true,
		UserRole:        "",
		SessionID:       "",
		IsAuthenticated: false,
		ErrorMessage:    "",
	}
}

//is authenticated

func (app *application) IsAuthenticated(r *http.Request) bool {
	session, err := app.sessionManager.Store.Get(r, "session")
	if err != nil {
		return false
	}
	userID, ok := session.Values["user_id"].(int)
	// Check if userID is present and is a non-zero value

	if !ok || userID == 0 {
		return false
	}
	return true
}
