package http

import (
	"net/http"
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
	h := NewHandler()
	mux.HandleFunc("/version", getVersion)
	mux.HandleFunc("/count", h.getCounter)
	mux.HandleFunc("/count_update", h.setCounter)
	return withTimer()(withCors(mux))
}
