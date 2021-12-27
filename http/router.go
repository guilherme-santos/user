package http

import (
	"net"
	"net/http"
	"time"

	"github.com/guilherme-santos/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(logger logrus.FieldLogger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if logger != nil {
		r.Use(Logger(logger))
	} else {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	return r
}

// Logger returns a request logging middleware
func Logger(log logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := user.SetLogger(r.Context(), log)
			reqID := middleware.GetReqID(ctx)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				elapsed := time.Since(t1)

				log := user.Logger(ctx)
				fields := logrus.Fields{
					"status_code": ww.Status(),
					"bytes":       ww.BytesWritten(),
					"elapsed_ms":  float64(elapsed.Nanoseconds()) / float64(time.Millisecond),
					"elapsed":     elapsed.String(),
					"remote_ip":   remoteIP,
					"proto":       r.Proto,
					"method":      r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				log.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
