// Package http provides the HTTP server for the application.
// we named this package `http` intentionally to force the
// isolation of the packages across this project
package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Run starts the HTTP server
func Run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	defer cancel()

	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)

	srv := NewServer()
	httpServer := &http.Server{
		Handler: srv.Handler(),
	}
	go func() {
		slog.Info("server listening...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Setup a ticker for periodic task
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				// make a new context for the Shutdown (thanks Alessandro Rosetti)
				shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				// do some cleanup works before exit
				if err := httpServer.Shutdown(shutdownCtx); err != nil {
					fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
				} else {
					slog.Info("shutting down server...")
				}
				return // returning to let the deferred cancel() run
			case <-ticker.C:
				// do something every minute
			}
		}
	}()
	wg.Wait()
	return nil
}
