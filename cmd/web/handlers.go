package main

import (
	"fmt"
	"html/template"

	"net/http"
	"strconv"

	"mayuraandrew.tech/snippetbox/pkg/models"
)

// define a home handler function which write a byte slice containing

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.

	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		app.notFound(w) // Use the notFound() helper
		return
	}
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// create an instance of a templateData struct holding the slice of snippets.

	data := &templateData{Snippets: s}

	files := []string {
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// pass in the templateData struct when executing the template.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}
// 	files := []string{
// 		"./ui/html/home.page.tmpl",
// 		"./ui/html/base.layout.tmpl",
// 		"./ui/html/footer.partial.tmpl",
// 	}
// // Initialize a slice containing the paths to the two files. Note that the
// // home.page.tmpl file must be the *first* file in the slice.

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		// app.errorLog.Println(err.Error())
// 		// http.Error(w, "Internal Server Error", 500)
// 		w.Write([]byte("Create a new snippet..."))
// 		app.serverError(w, err) // use the serverError()helper
// 		return
// 	}


// 	err = ts.Execute(w, nil)
// 		if err != nil {
// 			// app.errorLog.Println(err.Error())
// 			// http.Error(w, "Internal Server Error", 500)
// 			app.serverError(w, err) // use the serverError() helper
// 			return
// 		}


// add a showSnippet hanler function

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
// Extract the value of the id parameter from the query string and try to
// convert it to an integer using the strconv.Atoi() function. If it can't
// be converted to an integer, or the value is less than 1, we return a 404
// not found response.

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our respo
	// and write it to the http.ResponseWriter.
	
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partial that we made earlier.
	files := []string {
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// and then execute them. notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter

	// pass in the templateData strut when executing the template.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// add a createSnippet handler fucntion
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() method to write a "Method Not Allowed" response body. We
	// then return from the function so that the subsequent code is not execute
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
	// response header ma, errp. The first parameter is the header name, and
	// the second parameter is the header value.
		w.Header().Set("Allow", "POST")
		//w.WriteHeader(405)
		//http.Error(w, "Method Not Allowed", 405)
		//w.Write([]byte("Method Not Allowed"))
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// creaate some variables holding dummy data. We'll remove these later on during the build.
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi"
	expires := "7"

	// pass the data the SnippetModel.Insert() method, receiving the ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user tp the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}


