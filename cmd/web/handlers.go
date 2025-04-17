package main

import (
	"net/http"

	//	"github.com/cohune-cabbage/di/internal/data"
	//rangeable might use this
	"github.com/cohune-cabbage/di/internal/data"
	_ "golang.org/x/text/unicode/rangetable"
	//	"strconv"
	//	"github.com/cohune-cabbage/di/internal/validator"
)

// serverError writes an error message and stack trace to the error log, then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// handler for the homepage with error handling
// and template data

func (app *application) homepage(w http.ResponseWriter, r *http.Request) {

	
}
	