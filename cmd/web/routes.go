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
	mux.HandleFunc("/coach_dashboard", app.CoachDashBoardHandler)
	mux.HandleFunc("/interview", app.InterviewHandler)
	mux.HandleFunc("/loginuser", app.LoginHandler)
	return app.loggingMiddleware(mux)
}
