package main

import (
	"net/http"

	"github.com/cohune-cabbage/di/internal/data"
	"github.com/cohune-cabbage/di/internal/validator"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) homepage(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:           "Home",
		HeaderText:      "Welcome to V-Coach",
		PageDescription: "Your virtual coaching assistant.",
	}
	err := app.render(w, http.StatusOK, "homepage.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Sign Up",
		}
		err := app.render(w, http.StatusOK, "signup.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	signUp := &data.SignUp{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
		Role:  r.PostForm.Get("role"),
	}

	v := validator.NewValidator()
	v.Check(validator.NotBlank(signUp.Name), "name", "Name cannot be blank")
	v.Check(validator.NotBlank(signUp.Email), "email", "Email cannot be blank")
	v.Check(validator.IsValidEmail(signUp.Email), "email", "Invalid email address")

	if !v.ValidData() {
		data := &TemplateData{
			Title:        "Sign Up",
			SignUp:       signUp,
			SignUpErrors: v.Errors,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	// Save the signUp data to the database (omitted for brevity)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Login",
		}
		err := app.render(w, http.StatusOK, "login.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	login := &data.Login{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	v := validator.NewValidator()
	v.Check(validator.NotBlank(login.Email), "email", "Email cannot be blank")
	v.Check(validator.IsValidEmail(login.Email), "email", "Invalid email address")
	v.Check(validator.NotBlank(login.Password), "password", "Password cannot be blank")

	if !v.ValidData() {
		data := &TemplateData{
			Title:      "Login",
			Login:      login,
			FormErrors: v.Errors,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	// Authenticate the user (omitted for brevity)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) interviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &TemplateData{
			Title: "Interview",
			// Add other fields as needed
		}
		err := app.render(w, http.StatusOK, "interview.tmpl", data)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Load session, pull question from DB, and save response (omitted for brevity)

	http.Redirect(w, r, "/interview", http.StatusSeeOther)
}
