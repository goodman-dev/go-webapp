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

	log.Print("starting server on port 4000")

	// Start our web server
	// Omitting the host (host:port) means w elisten on all interfaces
	// If it returns an error, make sure we log it
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
