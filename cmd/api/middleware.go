package main

import (
	"net/http"
)

func (app *application) PermissionsValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := app.GetAuthStringFromHeader(w, r, "Authorization")
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, MissingAuthTokenError)
			return
		}
		app.logger.Info("Getting token", "token", tokenString)
		valid, err := app.services.TokenService.ValidateToken(tokenString, app.config.tokenConfig.secret)
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, MissingAuthTokenError)
			return
		}
		app.logger.Info("Validate token", "valid", valid)
		userId, err := app.ExtractUserIdFromToken(tokenString, app.config.tokenConfig.secret)
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}
		app.logger.Info("Getting user ID", "userId", userId)
		permissions, err := app.services.PermissionsService.GetPermissionsForUser(userId)
		if err != nil {
			app.serverSideErrorResponse(w, r, err)
			return
		}
		app.logger.Info("Validate token", "UsersPermissions", permissions)
		hasPermission := permissions.HasPermission("permissions:write")
		if !hasPermission {

			app.errorResponse(w, r, http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
