package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/oseau/web"
)

// Server is the http server
type Server struct{}

// NewServer creates a new server
func NewServer() Server {
	return Server{}
}

// Handler registers the routes and middlewares
// returns a ServeMux to be used by http.Server
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(web.VersionString))
	})
	var handler http.Handler = mux
	handler = withTimer()(handler)
	return handler
}

func withTimer() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				slog.Debug("req", "method", r.Method, "path", r.URL.Path, "time", time.Since(start).String())
			}(time.Now())
			handler.ServeHTTP(w, r)
		})
	}
}
