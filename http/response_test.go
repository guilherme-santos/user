package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guilherme-santos/user"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	testcases := []struct {
		Name       string
		Error      error
		StatusCode int
		Body       string
	}{
		{
			Name:       "InvalidArgument",
			Error:      user.NewMissingFieldError("name"),
			StatusCode: http.StatusBadRequest,
			Body:       `{"code": "missing_field", "field":"name", "message":"Field was not provided"}`,
		},
		{
			Name:       "NotFound",
			Error:      user.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Body:       `{"code": "not_found", "message":"user not found"}`,
		},
		{
			Name:       "Unknown",
			Error:      &user.Error{Type: user.Unknown, Code: "unknown_error", Message: "unknown error"},
			StatusCode: http.StatusInternalServerError,
			Body:       `{"code": "unknown_error", "message":"unknown error"}`,
		},
		{
			Name:       "GoError",
			Error:      errors.New("some other error"),
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"some other error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondWithError(w, tc.Error)
			assert.Equal(t, tc.StatusCode, w.Code)
			assert.JSONEq(t, tc.Body, w.Body.String())
		})
	}
}
