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

	return app.loggingMiddleware(mux)
}
