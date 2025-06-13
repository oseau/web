package http

import (
	"net/http"

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
	return withTimer()(withCors(mux))
}
