package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(app.routeNotFoundResponse)
	r.MethodNotAllowed(app.routeResourceNotAllowedResponse)

	r.Get("/healthcheck", app.healthcheckHandler)

	r.Post("/registerUser", app.registerUserHandler)
	r.Put("/registerUser", app.updateUserHandler)
	return r
}
