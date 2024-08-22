package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/goodman-dev/go-webapp/internal/models"

	// need the driver's init() function to run and register
	// the driver with the database/sql package
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr string
}

// Create an application struct for dependency injection
// If our handlers are spread across packages, use a closure pattern
// like https://gist.github.com/alexedwards/5cd712192b4831058b21
type application struct {
	logger   *slog.Logger
	cfg      *config
	snippets *models.SnippetModel
}

func main() {

	var cfg config

	// Parse cmd line flags into our config struct
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	// We need to use the parseTime=true parameter in our DSN to force it to convert
	// TIME and DATE fields to time.Time. Otherwise it returns these as []byte objects
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source")
	flag.Parse()

	// Set up our logger
	// ¡NOTE! The Logger created here is concurrency-safe, as long as we're
	// we're using the same slog.Logger instance for a destination
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	db, err := openDB(*dsn, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	defer db.Close()

	app := application{
		logger:   logger,
		cfg:      &cfg,
		snippets: &models.SnippetModel{DB: db},
	}

	// flag.String returns a pointer to a string, so dereference it
	logger.Info("starting server", slog.Group("request", slog.String("addr", cfg.addr)))

	// Start our web server
	err = http.ListenAndServe(cfg.addr, app.routes())
	// When app is terminated - or if we fail to start - log cause & exit
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string, logger *slog.Logger) (*sql.DB, error) {

	logger.Info("establishing DB connection pool", "dsn", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
