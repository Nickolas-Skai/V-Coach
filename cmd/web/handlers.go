package main

import (
	"net/http"

	//	"github.com/cohune-cabbage/di/internal/data"
	//rangeable might use this
	_ "golang.org/x/text/unicode/rangetable"
	//	"strconv"
	//	"github.com/cohune-cabbage/di/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data_homepage := &TemplateData{
		Title:           "Welcome",
		HeaderText:      "Welcome to the site",
		PageDescription: "This is the home page",
		NavLogo:         "/static/images/logo.svg",
	}
	err := app.render(w, http.StatusOK, "homepage.tmpl", data_homepage)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "homepage.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
