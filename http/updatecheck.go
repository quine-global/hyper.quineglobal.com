package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"app/html"
)

type updateCheckResponse struct {
	HasUpdate bool           `json:"has_update"`
	Stable    *releaseUpdate `json:"stable,omitempty"`
	Canary    *releaseUpdate `json:"canary,omitempty"`
}

type releaseUpdate struct {
	Version string `json:"version"`
	Build   int    `json:"build,omitempty"` // canary only
	Date    string `json:"date"`
	URL     string `json:"url,omitempty"`
}

// UpdateCheck registers GET /api/update-check.
//
// Query params:
//
//	version — current client version, e.g. "4.0.0" or "4.0.0-q-canary.12" (v-prefix stripped)
//	os      — mac | windows | linux | linux-rpm
//	arch    — x64 | arm64
//	canary  — "1" to also include canary update info (implied when client is already on canary)
func (s *Server) UpdateCheck(r chi.Router) {
	r.Get("/api/update-check", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		rawVersion := strings.TrimPrefix(q.Get("version"), "v")
		os := q.Get("os")
		arch := q.Get("arch")

		if rawVersion == "" || os == "" || arch == "" {
			http.Error(w, "version, os, and arch are required", http.StatusBadRequest)
			return
		}

		clientBase, clientBuild := parseTag(rawVersion)
		clientIsCanary := clientBuild > 0
		wantCanary := clientIsCanary || q.Get("canary") == "1"

		releases := s.releases.Get()

		var stableUpdate, canaryUpdate *releaseUpdate
		for _, rel := range releases {
			// Releases arrive newest-first from GitHub; take the first match on each track.
			if rel.IsCanary {
				if wantCanary && canaryUpdate == nil && newerCanary(rel, clientBase, clientBuild, clientIsCanary) {
					canaryUpdate = makeUpdate(rel, os, arch)
				}
			} else {
				if stableUpdate == nil && semverCmp(rel.Version, clientBase) > 0 {
					stableUpdate = makeUpdate(rel, os, arch)
				}
			}
			if stableUpdate != nil && (!wantCanary || canaryUpdate != nil) {
				break
			}
		}

		resp := updateCheckResponse{
			HasUpdate: stableUpdate != nil || canaryUpdate != nil,
			Stable:    stableUpdate,
			Canary:    canaryUpdate,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.log.Warn("Failed to encode update-check response", "error", err)
		}
	})
}

func makeUpdate(rel html.Release, os, arch string) *releaseUpdate {
	return &releaseUpdate{
		Version: rel.Version,
		Build:   rel.CanaryBuild,
		Date:    rel.Date,
		URL:     rel.DownloadURL(os, arch),
	}
}

// newerCanary reports whether rel is newer than the client's current canary position.
// For a canary client, the build number is the tie-breaker on the same base version.
// For a stable client opting into canary info, any canary on the same or higher base qualifies.
func newerCanary(rel html.Release, clientBase string, clientBuild int, clientIsCanary bool) bool {
	cmp := semverCmp(rel.Version, clientBase)
	if cmp > 0 {
		return true
	}
	if clientIsCanary {
		return cmp == 0 && rel.CanaryBuild > clientBuild
	}
	// Stable client: same-base canary is ahead of the stable release.
	return cmp == 0
}

// semverCmp compares two "major.minor.patch" version strings.
// Returns -1, 0, or 1.
func semverCmp(a, b string) int {
	pa := strings.SplitN(a, ".", 3)
	pb := strings.SplitN(b, ".", 3)
	for i := 0; i < 3; i++ {
		var na, nb int
		if i < len(pa) {
			na, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(pb[i])
		}
		if na != nb {
			if na < nb {
				return -1
			}
			return 1
		}
	}
	return 0
}
