package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// the serverError helper write an error message and stack trace to the errorlog
// then sends a generic 500 Internal Server Error reponse ot the user

func (app *application) serverError(w http.ResponseWriter, err error){

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}