package main

import (
	"authentication-service/internal/data"
	"authentication-service/internal/service"
	"net/http"
	"time"
)

func (app *application) RegenerateEmailTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input service.RegenerateEmailTokenInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	res, err := app.services.UserService.ValidateUser(input)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	if res.IsMatch == false {
		app.badRequestResponse(w, r, InvalidCombinationError)
		return
	}

	err = app.services.TokenService.DeleteTokensForUser(res.ID, data.ActivateEmailToken)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

	newToken, err := app.services.TokenService.CreateAccessToken(res.ID, data.ActivateEmailToken, app.config.tokenConfig.ttl, app.config.tokenConfig.secret)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	res.Token = newToken
	dataToken := &data.Token{
		Hash:   []byte(newToken),
		UserID: res.ID,
		Expiry: time.Now().Add(app.config.tokenConfig.ttl),
		Scope:  data.ActivateEmailToken,
	}
	app.logger.Info("Regenerate Email token", "Token", dataToken)
	_, err = app.services.TokenService.InsertToken(dataToken)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, res, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
}

func (app *application) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input service.ValidateTokenInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.Info("Validate token", "input", input)
	valid, err := app.services.TokenService.ValidateToken(input.Token, app.config.tokenConfig.secret)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	app.logger.Info("Validate token", "valid", valid)
	response := &service.ValidateTokenResponse{
		Token:   input.Token,
		IsValid: valid,
	}
	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

}
