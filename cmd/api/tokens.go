package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/denim-bluu/movie-db-app/internal/data"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var dt input
	err := app.readJSON(w, r, &dt)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByEmail(dt.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ok, err := user.Password.Matches(dt.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !ok {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, envelope{"authentication_token": token}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
