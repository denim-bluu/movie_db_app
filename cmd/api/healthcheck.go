package main

import (
	"net/http"
	"strconv"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment":     app.config.env,
			"version":         version,
			"limiter-enabled": strconv.FormatBool(app.config.limiter.enabled),
			"limiter-RPS":     strconv.FormatFloat(app.config.limiter.rps, 'f', -1, 64),
			"limiter-Burst":   strconv.Itoa(app.config.limiter.burst),
		},
	}
	err := app.writeJSON(w, env, http.StatusOK, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
}
