package main

import (
	"flag"
	"log/slog"
	"os"
)

type config struct {
	port int
	env  string

	jsonConfig struct {
		maxByte            int
		allowUnknownFields bool
	}
}
type application struct {
	config config
	logger *slog.Logger
}

const version = "1.0.0"

type responseData map[string]any

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&cfg.jsonConfig.maxByte, "Maxbytes", 1_048_576, "MAX Bytes for JSON Body")
	flag.BoolVar(&cfg.jsonConfig.allowUnknownFields, "AllowUnknownFields", false, "Allow unknown fields in JSON Body")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	logger.Error(err.Error())

}
