package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	res := responseData{
		"Status":      "Available",
		"Envirnoment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
	}
}
