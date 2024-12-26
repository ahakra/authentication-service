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

	r.Route("/v1", func(r chi.Router) {
		r.Post("/users", app.registerUserHandler)
		r.Put("/users", app.updateUserHandler)

		//check again register with domain password due to change in type

		r.Post("/auth/validateEmail", app.validateEmailHandler)
		r.Post("/auth/login", app.loginHandler)

		r.Post("/tokens/email", app.RegenerateEmailTokenHandler)
		r.Post("/tokens/validate", app.ValidateTokenHandler)

		r.Group(func(r chi.Router) {
			r.Use(app.PermissionsValidation)

			r.Post("/permissions", app.AddPermissionHandler)
			r.Post("/users/{userID}/permissions", app.AddPermissionToUserHandler)
			r.Delete("/users/{userID}/permissions", app.RemovePermissionFromUserHandler)
		})

	})

	return r
}
