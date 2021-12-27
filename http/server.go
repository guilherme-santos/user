package http

import (
	"net/http"
	"time"
)

// As a improvement we could pass a http.ServerConfig with more options
// and to avoid to have so many parameters.

func NewServer(addr string, router http.Handler) http.Server {
	return http.Server{
		Addr:              addr,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           router,
	}
}
