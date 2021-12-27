package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/guilherme-santos/user"
)

type UserHandler struct {
	svc user.Service
}

func NewUserHandler(r chi.Router, svc user.Service) *UserHandler {
	h := &UserHandler{
		svc: svc,
	}
	r.Post("/users", h.Create)
	r.Get("/users", h.List)
	r.Get("/users/{id}", h.Get)
	r.Put("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)
	return h
}

func (h UserHandler) Create(w http.ResponseWriter, req *http.Request) {
	var u *user.User

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		respondWithError(w, newJSONDecodeError(err))
		return
	}

	ctx := req.Context()

	err = h.svc.Create(ctx, u)
	if err != nil {
		respondWithError(w, err)
		return
	}

	u, err = h.svc.Get(ctx, u.ID)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondCreated(w, u)
}

func (h UserHandler) List(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	opts := user.NewListOptions()
	opts.Country = query.Get("country")
	opts.Sort = query.Get("sort")
	opts.Cursor = query.Get("cursor")
	if query.Has("per_page") {
		perPage, err := strconv.ParseInt(query.Get("per_page"), 10, 64)
		if err == nil {
			opts.PerPage = int64(perPage)
		}
	}
	if query.Has("page") {
		page, err := strconv.Atoi(query.Get("page"))
		if err == nil {
			opts.Page = page
		}
	}

	resp, err := h.svc.List(req.Context(), opts)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondOK(w, resp)
}

func (h UserHandler) Get(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	u, err := h.svc.Get(req.Context(), id)
	if err != nil {
		respondWithError(w, err)
		return
	}

	// I'd like to add some caching properties here, like the following
	// w.Header().Set("Cache-Control", "private, no-cache") // no-cache force the revalidation before use it
	// w.Header().Set("ETag", userEtag)

	// We could check If-None-Match and return the user only if changed
	// otherwise http.StatusNotModified

	respondOK(w, u)
}

func (h UserHandler) Update(w http.ResponseWriter, req *http.Request) {
	var u *user.User

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		respondWithError(w, newJSONDecodeError(err))
		return
	}

	u.ID = chi.URLParam(req, "id")
	ctx := req.Context()

	err = h.svc.Update(ctx, u)
	if err != nil {
		respondWithError(w, err)
		return
	}

	u, err = h.svc.Get(ctx, u.ID)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondOK(w, u)
}

func (h UserHandler) Delete(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	err := h.svc.Delete(req.Context(), id)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondNoContent(w)
}
