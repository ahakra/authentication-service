package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strings"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

	w.Header().Set("Content-Type", "application/json")
	res, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}
	res = append(res, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.WriteHeader(status)
	_, err = w.Write(res)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(app.config.jsonConfig.maxByte))

	// Create a JSON decoder
	dec := json.NewDecoder(r.Body)

	// Optionally disallow unknown fields
	if !app.config.jsonConfig.allowUnknownFields {
		dec.DisallowUnknownFields()
	}

	// Decode the JSON into the destination struct
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains incomplete JSON")

		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("body contains unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown field "))

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case err.Error() == "http: request body too large":
			return errors.New("body exceeds the maximum allowed size")

		default:
			return err // Return any other errors
		}
	}

	// Check for extra data after the main JSON value
	if dec.More() {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) GetAuthStringFromHeader(w http.ResponseWriter, r *http.Request, authKey string) (string, error) {

	authHeader := r.Header.Get(authKey)
	if authHeader == "" {

		return "", MissingAuthTokenError
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {

		return "", MissingAuthTokenError
	}
	return tokenString, nil
}

func (app *application) ExtractUserIdFromToken(tokenString string, secret string) (int64, error) {
	var secretKey = []byte(secret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("userID is missing or invalid in token")
	}

	return int64(userID), nil
}
