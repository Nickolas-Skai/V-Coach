package main

import (
	"net/http"
)

// routes sets up the application routes and handlers.

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.homepage)
	mux.HandleFunc("/login", app.LoginPageHandler)
	mux.HandleFunc("GET /signup", app.SignUpPageHandler)
	mux.HandleFunc("/user/signup", app.SignUpHandler)
	mux.Handle("/coach/dashboard", app.requireAuthentication(http.HandlerFunc(app.CoachDashboardHandler)))
	mux.Handle("/interview", app.requireAuthentication(http.HandlerFunc(app.InterviewHandler)))
	mux.HandleFunc("/loginuser", app.LoginHandler)
	mux.HandleFunc("/next/question", app.NextQuestionHandler)
	mux.HandleFunc("/previous/question", app.PreviousQuestionHandler)
	mux.Handle("/logout", app.requireAuthentication(http.HandlerFunc(app.LogoutHandler)))

	return app.loggingMiddleware(mux)
}
