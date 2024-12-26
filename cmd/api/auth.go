package main

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"authentication-service/internal/service"
	"fmt"
	"net/http"
	"time"
)

func (app *application) validateEmailHandler(w http.ResponseWriter, r *http.Request) {

	tokenString, err := app.GetAuthStringFromHeader(w, r, "Authorization")

	if err != nil {
		app.badRequestResponse(w, r, MissingVerificationTokenError)
		return
	}
	fmt.Printf("Validate Email token: %s\n", tokenString)

	// Validate the token
	validToken, err := app.services.TokenService.ValidateToken(tokenString, app.config.tokenConfig.secret)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, err.Error())
		_ = app.services.TokenService.DeleteToken([]byte(tokenString)) // Delete invalid token
		return
	}
	if !validToken {
		app.errorResponse(w, r, http.StatusUnauthorized, err.Error())
		_ = app.services.TokenService.DeleteToken([]byte(tokenString)) // Delete invalid token
		return
	}

	// Extract userId from the token
	userId, err := app.services.TokenService.ExtractUserIdFromToken(tokenString, app.config.tokenConfig.secret)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "Missing or Invalid Token")
		return
	}
	fmt.Printf("Validate Email userId from token: %d\n", userId)

	// Get the tokens for the user and validate if the token is correct
	tokens, err := app.services.TokenService.GetTokensForUserAndScope(userId, data.ActivateEmailToken)
	if err != nil || len(tokens) < 1 {
		app.errorResponse(w, r, http.StatusUnauthorized, "No activation token found")
		return
	}

	// Check if the valid token exists in the list of user tokens
	var validInput bool
	for _, token := range tokens {
		fmt.Printf("Validate Email tokens found for this user: %s\n", token.Hash)
		if string(token.Hash) == tokenString {
			validInput = true
			_ = app.services.TokenService.DeleteToken([]byte(tokenString)) // Delete the used token
			break
		}
	}

	// If the token is invalid, return error
	if !validInput {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid token")
		return
	}

	err = app.services.UserService.UpdateUserActivationStatus(userId, true)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	response := responseData{
		"data": "User email successfully validated and activated",
	}
	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var input service.LoginInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, opErr := app.services.UserService.GetUserByEmail(input.Email)
	if opErr != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

	pass := domain.Password{
		PasswordHash: user.Password,
	}
	match, err := pass.Matches(input.Password)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}

	if !match {
		app.badRequestResponse(w, r, InvalidCombinationError)
		return
	}
	err = app.services.TokenService.DeleteTokensForUser(user.ID, data.UserAccessToken)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	token, err := app.services.TokenService.CreateAccessToken(user.ID, data.UserAccessToken, app.config.tokenConfig.ttl, app.config.tokenConfig.secret)

	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	dataToken := &data.Token{
		Hash:   []byte(token),
		UserID: user.ID,
		Expiry: time.Now().Add(app.config.tokenConfig.ttl),
		Scope:  data.UserAccessToken,
	}
	_, err = app.services.TokenService.InsertToken(dataToken)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
	res := &service.LoginInResponse{
		AuthorizationToken: token,
	}
	err = app.writeJSON(w, http.StatusCreated, res, nil)
	if err != nil {
		app.serverSideErrorResponse(w, r, err)
		return
	}
}
