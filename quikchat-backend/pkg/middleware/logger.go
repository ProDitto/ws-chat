package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

// NewStructuredLogger is a middleware that logs requests using slog.
func NewStructuredLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				logger.Info("Request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"duration", time.Since(t1),
					"size", ww.BytesWritten(),
					"req_id", chi_middleware.GetReqID(r.Context()),
				)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

