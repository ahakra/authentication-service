package main

import (
	"authentication-service/internal/data"
	"authentication-service/internal/service"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input service.UserRegisterInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Info("User Input:", "value", input)
	userResponse, operationErrors := app.services.UserService.RegisterUser(&input)
	if operationErrors != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, operationErrors)
		return
	}

	jwt, err := app.services.TokenService.CreateAccessToken(userResponse.ID,
		data.ActivateEmailToken,
		app.config.tokenConfig.ttl,
		app.config.tokenConfig.secret)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err)
	}

	userResponse.VerificationToken = jwt
	err = app.writeJSON(w, http.StatusCreated, userResponse, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input service.UserRegisterInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	operationErrors := app.services.UserService.UpdateUser(&input)
	if operationErrors != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, operationErrors)
		return
	}

	app.writeJSON(w, http.StatusCreated, nil, nil)
}
