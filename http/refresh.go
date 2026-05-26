package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// InternalRefresh registers POST /internal/refresh to force an immediate GitHub release cache refresh.
func (s *Server) InternalRefresh(r chi.Router) {
	r.Post("/internal/refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := s.releases.ForceRefresh(r.Context()); err != nil {
			http.Error(w, "refresh failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "ok")
	})
}
