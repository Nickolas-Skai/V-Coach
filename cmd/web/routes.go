package main

//NEED UPDATE
import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.homepage)
	//	mux.HandleFunc("GET /login", app.loginpage)
	//	mux.HandleFunc("GET /signup", app.signuppage)

	return app.loggingMiddleware(mux)
}
