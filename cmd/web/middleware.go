package main

import (
	"net/http"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
		app.logger.Info("Request processed")
	})
	return fn
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.sessionManager.Exists(r, "user_id") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store, no-cache,")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := app.sessionManager.GetString(r, "user_role")
		if userRole != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
