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
	mux.HandleFunc("/interview/success", app.InterviewSuccessHandler)
	mux.Handle("/coach/teachers", app.requireRole("coach", http.HandlerFunc(app.GetAllTeachersHandler)))
	mux.HandleFunc("/coach/teachers//details", app.TeacherDetailHandler)
	mux.HandleFunc("/coach/teachers//sessions", app.AllInterviewSessionsListbyTeacherHandler)
	mux.HandleFunc("/interview_sessions/", app.InterviewDetailsHandler)
	mux.HandleFunc("/uploads", app.FileHandler)
	mux.Handle("/coach/sessions", app.requireRole("coach", http.HandlerFunc(app.AllSessionsHandler)))

	return app.loggingMiddleware(mux)
}
