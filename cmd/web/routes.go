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
	mux.HandleFunc("/signup", app.SignUpPageHandler)
	mux.HandleFunc("/User/signup", app.AddUserHandler)
	mux.HandleFunc("/coach_dashboard", app.CoachDashBoardHandler)
	mux.HandleFunc("/interview", app.InterviewHandler)
	mux.HandleFunc("/interview/submit", app.SubmitInterviewResponseHandler)

	return app.loggingMiddleware(mux)
}
