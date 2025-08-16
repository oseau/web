package http

import (
	"net/http"

	"github.com/rs/cors"
)

// URLFrontend is set at compile time using -ldflags
var URLFrontend = ""

func withCors(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{URLFrontend},
	}).Handler(h)
}
