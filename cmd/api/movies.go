package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/denim-bluu/movie-db-app/internal/data"
	"github.com/denim-bluu/movie-db-app/internal/validator"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Test",
		Runtime:   100,
		Genres:    []string{"test", "test2"},
		Version:   1,
	}

	v := validator.New()
	data.ValidateMovie(v, &movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, envelope{"movie": movie}, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}
