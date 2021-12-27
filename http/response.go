package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guilherme-santos/user"
)

var JSONContentType = "application/json; charset=UTF-8"

func respond(w http.ResponseWriter, status int, body interface{}) {
	if body != nil {
		w.Header().Set("Content-Type", JSONContentType)
	}
	w.WriteHeader(status)
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
}

func respondOK(w http.ResponseWriter, body interface{}) {
	respond(w, http.StatusOK, body)
}

func respondCreated(w http.ResponseWriter, body interface{}) {
	respond(w, http.StatusCreated, body)
}

func respondNoContent(w http.ResponseWriter) {
	respond(w, http.StatusNoContent, nil)
}

func respondWithError(w http.ResponseWriter, err error) {
	var uerr *user.Error
	if !errors.As(err, &uerr) {
		uerr = &user.Error{
			Type:    user.Unknown,
			Message: err.Error(),
		}
		// for this case we need to update err to be a valid json
		err = uerr
	}

	var status int

	switch uerr.Type {
	case user.InvalidArgument:
		status = http.StatusBadRequest
	case user.NotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	respond(w, status, err)
}

func newJSONDecodeError(err error) *user.Error {
	return &user.Error{
		Type:    user.InvalidArgument,
		Code:    "invalid_json",
		Message: err.Error(),
	}
}
