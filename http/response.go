package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guilherme-santos/user"
)

var JSONContentType = "application/json; charset=UTF-8"

func respondWithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", JSONContentType)

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

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
