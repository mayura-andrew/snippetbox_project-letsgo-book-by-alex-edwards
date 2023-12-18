package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)


func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		
		next.ServeHTTP(w, r)
	})
}


func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(){
			// using the buildin recover function to check if there has been a panic or not. if there was..
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				// call the app.serverError helper method to return a 500

				app.serverError(w, fmt.Errorf("%s", err))

			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent handlers in the chain are executed.
		if app.authenticatedUser(r) == 0 {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		
		// otherwise call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}


// create a NoSurf middleware function which uses a customized CSRF cooki with the Secure, Path and HttpOnly flags set.


func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})

	return csrfHandler
}