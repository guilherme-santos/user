package http

import (
	chilogrus "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(logger logrus.FieldLogger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if logger != nil {
		r.Use(chilogrus.Logger("router", logger))
	} else {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	return r
}
