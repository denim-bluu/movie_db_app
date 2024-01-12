package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (app *application) newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.NotFound(app.notFoundResponse)

	r.Use(app.metrics)
	r.Use(app.logRequest)
	r.Use(app.recoverPanic)
	r.Use(app.authenticate)
	r.Use(app.rateLimit)
	r.Use(app.enableCORS)

	// Require authentication
	r.Group(func(r chi.Router) {
		r.Use(app.requireReadPermission)
		r.Get("/v1/movies/{id}", app.showMovieHandler)
		r.Get("/v1/movies", app.listMoviesHandler)
		r.Post("/v1/movies", app.createMovieHandler)
		r.Patch("/v1/movies/{id}", app.updateMovieHandler)
		r.Delete("/v1/movies/{id}", app.deleteMovieHandler)
	})
	r.Group(func(r chi.Router) {
		r.Use(app.requireWritePermission)
		r.Post("/v1/movies", app.createMovieHandler)
		r.Patch("/v1/movies/{id}", app.updateMovieHandler)
		r.Delete("/v1/movies/{id}", app.deleteMovieHandler)
	})

	r.Group(func(r chi.Router) {
		r.Get("/debug/vars", expvar.Handler().ServeHTTP)
		r.Get("/v1/healthcheck", app.healthcheckHandler)
		r.Post("/v1/users", app.registerUserHandler)
		r.Put("/v1/users/activate", app.ActivateUserHandler)
		r.Post("/v1/tokens/authenticate", app.createAuthenticationTokenHandler)
	})
	return r
}
