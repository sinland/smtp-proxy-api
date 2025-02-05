package http

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler in the chain
		next.ServeHTTP(ww, r)

		// Log the request details
		duration := time.Since(start)
		slog.Info(
			fmt.Sprintf("%s %s [status=%d duration=%s]",
				r.Method,
				r.URL.Path,
				ww.Status(),
				duration,
			),
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", duration,
		)
	})
}
