package http

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/oseau/web/db"
	"github.com/oseau/web/redis"
	"github.com/oseau/web/ws"
	"golang.org/x/sync/errgroup"
)

// Server is the http server
type Server struct {
	handler *Handler
	db      *db.DB
	redis   *redis.Redis
	hub     *ws.Hub
}

var (
	onceServer sync.Once
	srv        *Server
)

// NewServer creates a new server
func NewServer() *Server {
	onceServer.Do(func() {
		srv = &Server{
			db:    db.NewDB(),
			redis: redis.NewRedis(),
			hub:   ws.NewHub(),
		}
		go srv.hub.Run()
		srv.handler = NewHandler(srv.db, srv.redis, srv.hub)
	})
	return srv
}

// Handler registers the routes and middlewares
// returns a ServeMux to be used by http.Server
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", getVersion)
	mux.HandleFunc("/count", s.handler.getCounter)
	mux.HandleFunc("/count-click", s.handler.setCounter)
	mux.HandleFunc("/ws", s.handler.ws)
	return withTimer()(withCors(mux))
}

// Close closes the server, releasing db, redis and websocket resources
func (s *Server) Close(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		done := make(chan error, 1)
		go func() {
			done <- s.db.Close()
		}()
		select {
		case err := <-done:
			if err != nil {
				slog.Error("failed to close db", "error", err)
			}
			slog.Info("database closed successfully")
			return err
		case <-ctx.Done():
			slog.Error("timeout closing database")
			return ctx.Err()
		}
	})
	g.Go(func() error {
		done := make(chan error, 1)
		go func() {
			done <- s.redis.Close()
		}()
		select {
		case err := <-done:
			if err != nil {
				slog.Error("failed to close redis", "error", err)
			}
			slog.Info("redis closed successfully")
			return err
		case <-ctx.Done():
			slog.Error("timeout closing redis")
			return ctx.Err()
		}
	})
	g.Go(func() error {
		done := make(chan error, 1)
		go func() {
			done <- s.hub.Close()
		}()
		select {
		case err := <-done:
			if err != nil {
				slog.Error("failed to close websocket", "error", err)
			}
			slog.Info("websocket closed successfully")
			return err
		case <-ctx.Done():
			slog.Error("timeout closing websocket")
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		slog.Error("some resources failed to close", "error", err)
		return err
	}
	slog.Info("all resources closed successfully")
	return nil
}
