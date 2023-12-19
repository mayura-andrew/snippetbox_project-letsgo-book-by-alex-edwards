package main

import (
	"fmt"
	"net/http"
	"strconv"

	"mayuraandrew.tech/snippetbox/pkg/forms"
	"mayuraandrew.tech/snippetbox/pkg/models"
)

// define a home handler function which write a byte slice containing

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.

	// if r.URL.Path != "/" {
	// 	//http.NotFound(w, r)
	// 	app.notFound(w) // Use the notFound() helper
	// 	return
	// }

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}


	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

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

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

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

	//flash := app.session.PopString(r, "flash")

	// use the new render helper 
	app.render(w, r, "show.page.tmpl", &templateData{
	//	Flash: flash,
		Snippet: s,
	})
}

// add a createSnippet handler fucntion
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() method to write a "Method Not Allowed" response body. We
	// then return from the function so that the subsequent code is not execute
	// if r.Method != "POST" {
	// 	// Use the Header().Set() method to add an 'Allow: POST' header to the
	// // response header ma, errp. The first parameter is the header name, and
	// // the second parameter is the header value.
	// 	w.Header().Set("Allow", "POST")
	// 	//w.WriteHeader(405)
	// 	//http.Error(w, "Method Not Allowed", 405)
	// 	//w.Write([]byte("Method Not Allowed"))
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires := r.PostForm.Get("expires")

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	// create a new snippet record in the database using the form data.


	// errors := make(map[string]string)

	// if strings.TrimSpace(title) == ""{
	// 	errors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	errors["title"] = "This filed is too long (maximum is 100 characters)"
	// }

	// if strings.TrimSpace(content) == "" {
	// 	errors["content"] = "This field cannit be blank"
	// }

	// if strings.TrimSpace(expires) == "" {
	// 	errors["expires"] = "This field cannot be blank"
	// } else if expires != "365" && expires != "7" && expires != "1" {
	// 	errors["expires"] = "This field is invalid"
	// }

	// if len(errors) > 0 {
	// 	app.render(w, r, "create.page.tmpl", &templateData{
	// 		FormErrors: errors,
	// 		FormData: r.PostForm,
	// 	})

	// 	return
	// }

	if !form.Valid(){
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")
    expires := r.PostForm.Get("expires")

    app.infoLog.Println("Form values:", title, content, expires)

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	// redirect the user tp the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}



func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request){
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}


func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validate the form contents using the form helper we made earlier.

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid(){
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
	}

	// try to create a new user record in the database
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Addrees is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// otherwise add a confirmation flash message to the session confirming that their signup worked and asking them to log in

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// and redirect the user to the login page.

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	// otherwise send a placeholder respone 
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request){
	app.session.Remove(r, "userID")

	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", 303)
}