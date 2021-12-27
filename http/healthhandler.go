package http

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(r chi.Router, db *sql.DB) *HealthHandler {
	h := &HealthHandler{
		db: db,
	}
	r.Get("/health", h.Health)
	return h
}

func (h HealthHandler) Health(w http.ResponseWriter, req *http.Request) {
	res := map[string]interface{}{
		"status": "ok",
	}
	dbres := map[string]interface{}{}

	err := h.db.PingContext(req.Context())
	if err != nil {
		dbres["status"] = "error"
		dbres["error"] = err.Error()
	} else {
		dbres["status"] = "ok"
		stats := h.db.Stats()
		dbres["open_conns"] = stats.OpenConnections
		dbres["in_use_conns"] = stats.InUse
		dbres["idle_conns"] = stats.Idle
	}

	res["mysql"] = dbres
	respondOK(w, res)
}
