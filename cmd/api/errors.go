package main

import (
	"errors"
	"fmt"
	"net/http"
)

var InvalidCombinationError = errors.New("invalid combination")

func (app *application) logError(r *http.Request, err error) {
	var method = r.Method
	var uri = r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	res := responseData{"error": data}
	err := app.writeJSON(w, status, res, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverSideErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logError(r, err)
	message := "internal server error"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) routeNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "route not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) routeResourceNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the method %s is not supported for this route", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
