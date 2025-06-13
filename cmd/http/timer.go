package http

import (
	"log/slog"
	"net/http"
	"time"
)

func withTimer() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				took := time.Since(start)
				if took > time.Second {
					slog.Info("req", "method", r.Method, "path", r.URL.Path, "time", took.String())
				} else {
					slog.Debug("req", "method", r.Method, "path", r.URL.Path, "time", took.String())
				}
			}(time.Now())
			handler.ServeHTTP(w, r)
		})
	}
}
