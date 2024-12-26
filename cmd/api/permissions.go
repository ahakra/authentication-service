package main

import (
	"authentication-service/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) AddPermissionHandler(w http.ResponseWriter, r *http.Request) {

	var input service.AddPermissionInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.PermissionsService.AddPermission(input.Permission)
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

func (app *application) AddPermissionToUserHandler(w http.ResponseWriter, r *http.Request) {

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64) // Convert the string to int64
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var input service.AddPermissionInput
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.services.PermissionsService.AddPermissionToUser(userID, input.Permission)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, nil, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

}

func (app *application) RemovePermissionFromUserHandler(w http.ResponseWriter, r *http.Request) {

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64) // Convert the string to int64
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var input service.DeletePermissionInput
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.PermissionsService.RemovePermission(userID, input.Permission)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, nil, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

}
