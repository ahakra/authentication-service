package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *application) serve() error {

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.config.port),
		Handler:        app.routes(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("Server starting", "port", app.config.port)
	err := s.ListenAndServe()
	return err
}
