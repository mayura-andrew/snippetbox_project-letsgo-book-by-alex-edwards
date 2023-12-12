package main

import (
	"flag"
	"log"
	"net/http"
	"os"

)


// define an application struct to hold the application wide dependencies for the 
// web application. For now we'll only include fields for the two custom logger 
// we'll add more to it as the build progresses.

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}
func main() {
	// Define a new command-line flag with the name 'addr', a default value of
	// and some short help text explaining what the flag controls. The value of
	// flag will be stored in the addr variable at runtime.

	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any
	// encountered during parsing the application will be terminated.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)


	// initialize a new instance of application containing the dependencies.

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
	// use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)
	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.

	// Create a file server which serves files out of the "./ui/static" directo
	// Note that the path given to the http.Dir function is relative to the pro
	// directory root.
//fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with "/static/". For matching paths, we strip t
	// "/static" prefix before the request reaches the file server

//	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(), // call the new app.routes() method
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
