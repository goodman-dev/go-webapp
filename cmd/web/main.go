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

// Create an application struct for dependency injection
// If our handlers are spread across packages, use a closure pattern
// like https://gist.github.com/alexedwards/5cd712192b4831058b21
type application struct {
	logger *slog.Logger
	cfg    *config
}

func main() {

	var cfg config

	// Parse cmd line flags into our config struct
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	// Set up our logger
	// ¡NOTE! The Logger created here is concurrency-safe, as long as we're
	// we're using the same slog.Logger instance for a destination
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	app := application{
		logger: logger,
		cfg:    &cfg,
	}

	// flag.String returns a pointer to a string, so dereference it
	logger.Info("starting server", slog.Group("request", slog.String("addr", cfg.addr)))

	// Start our web server
	err := http.ListenAndServe(cfg.addr, app.routes())
	// When app is terminated - or if we fail to start - log cause & exit
	logger.Error(err.Error())
	os.Exit(1)

}
