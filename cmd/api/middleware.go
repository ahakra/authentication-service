package main

import "net/http"

func (app *application) PermissionsValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := app.GetAuthStringFromHeader(w, r, "authorization_token")
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, MissingAuthTokenError)
			return
		}
		valid, err := app.services.TokenService.ValidateToken(tokenString, app.config.tokenConfig.secret)
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, MissingAuthTokenError)
			return
		}
		app.logger.Info("Validate token", "valid", valid)
		userId, err := app.services.TokenService.ExtractUserIdFromToken(tokenString, app.config.tokenConfig.secret)
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}
		permissions, err := app.services.PermissionsService.GetPermissionsForUser(userId)
		if err != nil {
			app.serverSideErrorResponse(w, r, err)
			return
		}
		hasPermission := permissions.HasPermission("permissions:write")
		if !hasPermission {
			app.errorResponse(w, r, http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
