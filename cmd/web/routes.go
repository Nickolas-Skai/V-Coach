package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.homepage)
	mux.HandleFunc("/login", app.LoginPageHandler)
	mux.HandleFunc("GET /signup", app.SignUpPageHandler)
	mux.HandleFunc("/user/signup", app.SignUpHandler)
	mux.Handle("/coach_dashboard", app.requireAuthentication(http.HandlerFunc(app.CoachDashBoardHandler)))
	mux.Handle("/interview", app.requireAuthentication(http.HandlerFunc(app.InterviewHandler)))
	mux.HandleFunc("/loginuser", app.LoginHandler)
	mux.Handle("/logout", app.requireAuthentication(http.HandlerFunc(app.LogoutHandler)))
	return app.loggingMiddleware(mux)
}
