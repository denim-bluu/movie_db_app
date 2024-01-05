package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (app *application) newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.MethodNotAllowed(app.methodNotAllowedResponse)
	router.NotFound(app.notFoundResponse)

	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Post("/v1/movies", app.createMovieHandler)
	router.Get("/v1/movies/{id}", app.showMovieHandler)
	router.Patch("/v1/movies/{id}", app.updateMovieHandler)
	router.Delete("/v1/movies/{id}", app.deleteMovieHandler)
	router.Get("/v1/movies", app.listMoviesHandler)
	return router
}
