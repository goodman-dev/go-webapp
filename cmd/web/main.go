package main

import (
	"log"
	"net/http"
)

func main() {

	// Create a new http request multiplexer & register a route
	mux := http.NewServeMux()

	// ¡NOTE! servemux treates the route pattern "/" like a catch-all
	// In fact, this applies to all trailing-slash paths
	// So we can use regex sequence {$} to prevent this
	// When paths match multiple routes, "the most specific route wins"

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreateForm)
	mux.HandleFunc("POST /snippet/create", snippetCreate)

	// Create a file handler that serves from the ./ui/static dir,
	// and strip /static/ from the start of the URL path so that e.g.,
	// <host>/static/css/main.css reaches ./ui/static/css/main.css
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	log.Print("starting server on port 4000")

	// Start our web server
	// Omitting the host (host:port) means w elisten on all interfaces
	// If it returns an error, make sure we log it
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
