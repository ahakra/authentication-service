package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	w.Write(res)

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
