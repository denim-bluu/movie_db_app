package main

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/denim-bluu/movie-db-app/internal/data"
	"github.com/denim-bluu/movie-db-app/internal/mailer"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	Smtp mailer.Smtp
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models *data.Models
	mailer *mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := parseConfig()
	handleError(err, logger)

	db, err := openDB(cfg)
	handleError(err, logger)
	defer db.Close()

	InitExpvar(db)

	app := createApplication(cfg, logger, db)

	err = app.serve()
	handleError(err, logger)
}
