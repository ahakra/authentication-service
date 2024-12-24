package main

import (
	"authentication-service/internal/data"
	"authentication-service/internal/service"
	"context"
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"time"
)

type config struct {
	port int
	env  string

	jsonConfig struct {
		maxByte            int
		allowUnknownFields bool
	}
	tokenConfig struct {
		secret string
		ttl    time.Duration
	}

	db struct {
		dsn string
	}
}
type application struct {
	config   config
	logger   *slog.Logger
	services *service.ServiceManager
}

const version = "1.0.0"

type responseData map[string]any

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&cfg.jsonConfig.maxByte, "Maxbytes", 1_048_576, "MAX Bytes for JSON Body")
	flag.BoolVar(&cfg.jsonConfig.allowUnknownFields, "AllowUnknownFields", false, "Allow unknown fields in JSON Body")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "./database.db", "SQLITE3 DSN")

	flag.StringVar(&cfg.tokenConfig.secret, "secret", "defaultSecret", "The secret key for token signing")
	flag.DurationVar(&cfg.tokenConfig.ttl, "ttl", 3*24*time.Hour, "The time-to-live for the token")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	userRepo := data.NewUserRepository(db)
	tokenRepo := data.NewTokenRepository(db)
	permissionsRepo := data.NewPermissionsRepository(db)
	repoManager := data.NewRepoManager(userRepo, tokenRepo, permissionsRepo)

	userService := service.NewUserService(repoManager)
	tokenService := service.NewTokenService(repoManager)
	permissionsService := service.NewPermissionsService(repoManager)

	serviceManager := service.NewServiceManager(userService, tokenService, permissionsService)
	app.services = serviceManager

	err = app.serve()
	logger.Error(err.Error())

}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
