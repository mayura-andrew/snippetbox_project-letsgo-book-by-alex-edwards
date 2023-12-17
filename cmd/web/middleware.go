package main

import (
	"fmt"
	"net/http"
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