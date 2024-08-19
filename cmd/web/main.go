package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr string
}

func main() {

	var cfg config

	// Parse cmd line flags into our config struct
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	// Create a new http request multiplexer & register a route
	mux := http.NewServeMux()

	// Set up our logger
	// ¡NOTE! The Logger created here is concurrency-safe, as long as we're
	// we're using the same slog.Logger instance for a destination
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

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

	// flag.String returns a pointer to a string, so dereference it
	logger.Info("starting server", slog.Group("request", slog.String("addr", cfg.addr)))

	// Start our web server
	err := http.ListenAndServe(cfg.addr, mux)
	// When app is terminated - or if we fail to start - log cause & exit
	logger.Error(err.Error())
	os.Exit(1)

}
