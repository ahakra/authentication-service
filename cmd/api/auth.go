package main

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"authentication-service/internal/service"
	"net/http"
	"strings"
	"time"
)

func (app *application) validateEmailHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.errorResponse(w, r, http.StatusUnauthorized, "Missing Authorization Token")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		app.errorResponse(w, r, http.StatusUnauthorized, "Missing Authorization Token")
		return
	}

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

	// Get the tokens for the user and validate if the token is correct
	tokens, err := app.services.TokenService.GetTokensForUserAndScope(userId, data.ActivateEmailToken)
	if err != nil || len(tokens) < 1 {
		app.errorResponse(w, r, http.StatusUnauthorized, "No activation token found")
		return
	}

	// Check if the valid token exists in the list of user tokens
	var validInput bool
	for _, token := range tokens {
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

	// Retrieve user by ID
	user, opp := app.services.UserService.GetUserByID(userId)
	if opp != nil && (opp.Validation != nil || opp.Database != nil) {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, opp)
		return
	}

	// Update user to set isActivated to true
	user.Activated = true
	userUpdateInput := &service.UserRegisterInput{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Activated: user.Activated,
	}
	opp = app.services.UserService.UpdateUser(userUpdateInput)
	if opp != nil && (opp.Validation != nil || opp.Database != nil) {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, opp)
		return
	}

	// Return a success response
	resp := responseData{
		"data": "User email successfully validated and activated",
	}
	err = app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var input service.LoginInput
	err := app.readJSON(w, r, input)
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
