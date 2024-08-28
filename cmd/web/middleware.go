package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware function that logs information about incoming HTTP requests.
// It takes an http.Handler as input and returns an http.Handler.
// The returned http.Handler logs the remote address, protocol, HTTP method, and URL path of each request.
// The logged information is written to the application's infoLog.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware function that recovers from panics in the application.
// It wraps the provided http.Handler and recovers from any panics that occur during its execution.
// If a panic occurs, it sets the "Connection" header to "close" and calls the serverError method of the application,
// passing the recovered error as a parameter.
// This middleware ensures that the application continues to run even if a panic occurs, preventing the server from crashing.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requireAuthentication is a middleware function that checks if the user is authenticated before allowing access to the next handler.
// If the user is not authenticated, it redirects them to the login page.
// It also adds a "Cache-Control" header with the value "no-store" to the response.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// authenticate is a middleware function that checks if a user is authenticated.
// It takes the next http.Handler as a parameter and returns an http.Handler.
// If the user is not authenticated, the next handler is called.
// If the user is authenticated, it checks if the user exists in the database.
// If the user exists, it adds a context value to the request indicating that the user is authenticated.
// Finally, it calls the next handler.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// noSurf is a middleware function that adds CSRF protection to the HTTP handler chain.
// It wraps the provided handler with a nosurf.CSRFHandler and sets the base cookie with the following properties:
// - HttpOnly: true
// - Path: "/"
// - Secure: true
// The CSRFHandler ensures that all incoming requests have a valid CSRF token, preventing cross-site request forgery attacks.
// It returns the CSRFHandler as an http.Handler.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
