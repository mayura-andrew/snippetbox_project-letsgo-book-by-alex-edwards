package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"mayuraandrew.tech/snippetbox/pkg/models"
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
		if app.authenticatedUser(r) == nil {
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if a userID value exists in the session. If this isn't present then call next handler in the chain as normal.
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// fetch the details of the current user from the database.
		// if no matchinng record is found, remove the (invalid) userID from
		// their session and call the next handler in the chain as normal.
		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		}  else if err != nil {
			app.serverError(w, err)
			return
		}

		// otherwise, we know that the request is coming from a valid.
		// authenticated (logged in) user. We create a new copy of the 
		// requesst with the user information added to the request context, and
		// call the next handler in the chain *using this new copy of the request.
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}