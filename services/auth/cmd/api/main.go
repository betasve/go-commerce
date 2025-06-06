package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/data"
	"github.com/betasve/go-commerce/services/auth/internal/jsonlog"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn     string
		migrate string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "development|staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GO_COMMERCE_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.db.migrate, "db-migrate", "false", "Trigger DB Migration")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if cfg.db.migrate == "true" {
		err := migrateDB(db)

		if err != nil {
			logger.PrintFatal(err, nil)
		} else {
			os.Exit(0)
		}
	}

	logger.PrintInfo("starting server", map[string]string{
		"env":  cfg.env,
		"port": strconv.Itoa(cfg.port),
	})

	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sql.DB) error {
	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s/db/migrations/", os.Getenv("SERVICE_ROOT_PATH")),
		"postgres",
		migrationDriver,
	)

	if err != nil {
		return err
	}

	err = migrator.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
