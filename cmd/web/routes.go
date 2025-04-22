package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.homepage)
	mux.HandleFunc("/login", app.loginHandler)
	mux.HandleFunc("/signup", app.signupHandler)
	mux.HandleFunc("/interview", app.InterviewHandler)
	mux.HandleFunc("/coach/dashboard", app.CoachDashboardHandler)
	mux.HandleFunc("/questions/manage", app.ManageQuestionsHandler)
	mux.HandleFunc("/questions/edit/{id}", app.EditQuestionHandler)
	mux.HandleFunc("/questions/new", app.NewQuestionHandler)
	mux.HandleFunc("/coach/sessions/{teacher_id}", app.TeacherSessionsHandler)
	mux.HandleFunc("/coach/delete_teacher/{id}", app.DeleteTeacherHandler)

	return app.loggingMiddleware(mux)
}
