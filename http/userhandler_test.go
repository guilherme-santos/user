package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/guilherme-santos/user"
	uhttp "github.com/guilherme-santos/user/http"
	"github.com/guilherme-santos/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandlerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := newUser()
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(u)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", reqBody)

	// Creates the user
	svc := mock.NewUserService(ctrl)
	svc.EXPECT().
		Create(gomock.Any(), u).
		DoAndReturn(func(_ context.Context, u *user.User) error {
			u.ID = "uuid"
			return nil
		})
	svc.EXPECT().Get(gomock.Any(), "uuid").Return(u, nil)

	r := uhttp.NewRouter()
	uhttp.NewUserHandler(r, svc)
	r.ServeHTTP(w, req)

	userJSON, _ := json.Marshal(u)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, string(userJSON), w.Body.String())
}

func TestUserHandlerList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u1 := newUser()
	u1.ID = "uuid-1"

	u2 := newUser()
	u2.ID = "uuid-2"

	u3 := newUser()
	u3.ID = "uuid-3"

	w := httptest.NewRecorder()
	query := url.Values{}
	query.Set("country", "DE")
	query.Set("sort", "first_name")
	query.Set("cursor", "abc")
	query.Set("per_page", "50")
	query.Set("page", "2")

	req, _ := http.NewRequest(http.MethodGet, "/users?"+query.Encode(), nil)

	resp := &user.ListResponse{
		Total:      3,
		PerPage:    50,
		Users:      []*user.User{u1, u2, u3},
		NextCursor: "def",
	}

	// Retrieves list of user
	svc := mock.NewUserService(ctrl)
	svc.EXPECT().
		List(gomock.Any(), &user.ListOptions{
			Country: "DE",
			Sort:    "first_name",
			Cursor:  "abc",
			PerPage: 50,
			Page:    2,
		}).
		Return(resp, nil)

	r := uhttp.NewRouter()
	uhttp.NewUserHandler(r, svc)
	r.ServeHTTP(w, req)

	respJSON, _ := json.Marshal(resp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(respJSON), w.Body.String())
}

func TestUserHandlerGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := newUser()
	u.ID = "uuid"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/uuid", nil)

	// Retrieves the user
	svc := mock.NewUserService(ctrl)
	svc.EXPECT().Get(gomock.Any(), "uuid").Return(u, nil)

	r := uhttp.NewRouter()
	uhttp.NewUserHandler(r, svc)
	r.ServeHTTP(w, req)

	userJSON, _ := json.Marshal(u)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(userJSON), w.Body.String())
}

func TestUserHandlerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := newUser()
	u.ID = "uuid"
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(u)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/uuid", reqBody)

	// Updates the user
	svc := mock.NewUserService(ctrl)
	svc.EXPECT().Update(gomock.Any(), u).Return(nil)
	svc.EXPECT().Get(gomock.Any(), "uuid").Return(u, nil)

	r := uhttp.NewRouter()
	uhttp.NewUserHandler(r, svc)
	r.ServeHTTP(w, req)

	userJSON, _ := json.Marshal(u)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(userJSON), w.Body.String())
}

func TestUserHandlerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/uuid", nil)

	// Deletes the user
	svc := mock.NewUserService(ctrl)
	svc.EXPECT().Delete(gomock.Any(), "uuid").Return(nil)

	r := uhttp.NewRouter()
	uhttp.NewUserHandler(r, svc)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func newUser() *user.User {
	return &user.User{
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}
}
