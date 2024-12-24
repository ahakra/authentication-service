package main

import (
	"authentication-service/internal/service"
	"net/http"
)

func (app *application) AddPermissionHandler(w http.ResponseWriter, r *http.Request) {

	var input service.AddPermissionInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.PermissionsService.AddPermission(input.UserID, input.Permission)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, input, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

}

func (app *application) RemovePermissionHandler(w http.ResponseWriter, r *http.Request) {

	var input service.RemovePermissionInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.PermissionsService.RemovePermission(input.UserID, input.Permission)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, input, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

}
