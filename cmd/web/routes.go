package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	// Create a new http request multiplexer & register our routes
	mux := http.NewServeMux()

	// ¡NOTE! servemux treates the route pattern "/" like a catch-all
	// In fact, this applies to all trailing-slash paths
	// So we can use regex sequence {$} to prevent this
	// When paths match multiple routes, "the most specific route wins"

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreateForm)
	mux.HandleFunc("POST /snippet/create", app.snippetCreate)

	// Create a file handler that serves from the ./ui/static dir,
	// and strip /static/ from the start of the URL path so that e.g.,
	// <host>/static/css/main.css reaches ./ui/static/css/main.css
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	return mux

}
