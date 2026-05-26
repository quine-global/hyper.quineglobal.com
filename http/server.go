// Package http has the [Server] and HTTP handlers.
package http

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mrz1836/postmark"
)

// Server holds dependencies for the HTTP server as well as the HTTP server itself.
type Server struct {
	log      *slog.Logger
	mux      chi.Router
	server   *http.Server
	postmark *postmark.Client
	releases *ReleaseCache
}

type NewServerOptions struct {
	Log           *slog.Logger
	PostmarkToken string
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	mux := chi.NewMux()

	return &Server{
		log:      opts.Log,
		mux:      mux,
		postmark: postmark.NewClient(opts.PostmarkToken, ""),
		releases: &ReleaseCache{log: opts.Log},
		server: &http.Server{
			Addr:              ":8081",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the server and set up routes.
func (s *Server) Start() error {
	s.log.Info("Starting http server", "address", "http://localhost:8081")

	s.setupRoutes()

	// Fetch releases immediately, then refresh every hour in the background.
	go s.releases.TryRefresh(context.Background(), 0)
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			s.releases.TryRefresh(context.Background(), time.Hour)
		}
	}()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the server gracefully.
func (s *Server) Stop() error {
	s.log.Info("Stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Info("Stopped http server")
	return nil
}
