package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the User.
func (app *application) serverError(w http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2,trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
}

// The ClientError helper sends a specific status code and cooresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int){
	http.Error(w,http.StatusText(status),status)
}

// Not Found Helper
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w,http.StatusNotFound)
}