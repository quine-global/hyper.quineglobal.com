package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"app/html"
)

// hazelUpdate is the response format expected by Electron autoUpdater
// (the same protocol used by the original releases.hyper.is / Hazel).
type hazelUpdate struct {
	Name    string `json:"name"`
	Notes   string `json:"notes,omitempty"`
	PubDate string `json:"pub_date,omitempty"`
	URL     string `json:"url,omitempty"`
}

// UpdateCheck registers the auto-update feed endpoints.
//
//   - GET /api/update-check  (primary, matches existing expectations)
//   - GET /update            (convenience short path)
//
// This implements the update protocol used by Hyper's Electron autoUpdater
// (the same one originally served from releases.hyper.is via Hazel).
//
// Query params:
//
//	version  — client's current version, e.g. "4.0.0-q-canary.12" or "v4.0.0"
//	platform — darwin | win32 | linux | deb   (or os= for compatibility)
//	arch     — x64 | arm64
//	channel  — "stable" | "canary"   (explicit preference; takes precedence for track selection)
//
// Response:
//   - 204 No Content  → no update available
//   - 200 JSON        → update available (Hazel/Squirrel format: {name, notes, pub_date, url?})
func (s *Server) UpdateCheck(r chi.Router) {
	r.Get("/api/update-check", s.updateCheckHandler())
}

func (s *Server) updateCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		rawVersion := strings.TrimPrefix(q.Get("version"), "v")
		platform := firstNonEmpty(q.Get("platform"), q.Get("os"))
		arch := q.Get("arch")
		channel := strings.ToLower(q.Get("channel"))

		if rawVersion == "" {
			http.Error(w, "version is required", http.StatusBadRequest)
			return
		}

		clientBase, clientBuild := parseTag(rawVersion)
		clientIsCanary := clientBuild > 0 || channel == "canary" || channel == "canaryupdates"

		// Map incoming platform names from the Hyper client to our internal asset keys.
		internalOS := normalizePlatformForAssets(platform)

		releases := s.releases.Get()

		// Find the newest release that is newer than what the client has,
		// on the correct track (stable vs canary), and preferably has assets
		// for the requested platform/arch.
		var candidate *html.Release
		for _, rel := range releases {
			if rel.IsCanary {
				if clientIsCanary && newerCanary(rel, clientBase, clientBuild, true) {
					candidate = &rel
					break
				}
				// Non-canary client can still be offered canary if they opted in,
				// but by default we only auto-offer canary to canary clients
				// (matching original behavior of separate canary/stable feeds).
				if !clientIsCanary && newerCanary(rel, clientBase, clientBuild, false) {
					candidate = &rel
					break
				}
			} else {
				if !clientIsCanary && semverCmp(rel.Version, clientBase) > 0 {
					candidate = &rel
					break
				}
			}
		}

		if candidate == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Build the Hazel-compatible response.
		resp := hazelUpdate{
			Name:    candidate.TagName,
			Notes:   fmt.Sprintf("Quine Hyper %s", candidate.TagName),
			PubDate: candidate.PublishedAt.UTC().Format(time.RFC3339),
			URL:     candidate.DownloadURL(internalOS, arch),
		}

		// If we have no precise asset URL for this platform/arch, still tell the
		// client there is an update (especially important for Linux where the
		// client only uses the feed to learn *that* a new version exists).
		if resp.URL == "" {
			// Provide a generic GitHub releases URL as a fallback.
			// The client already falls back to constructing a tag URL.
			resp.URL = fmt.Sprintf("https://github.com/quine-global/hyper/releases/tag/%s", candidate.TagName)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.log.Warn("Failed to encode update response", "error", err)
		}
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// normalizePlatformForAssets converts values the Hyper client sends
// (darwin, win32, deb, etc.) into the keys used in Release.Assets.
func normalizePlatformForAssets(p string) string {
	switch strings.ToLower(p) {
	case "darwin", "mac", "osx":
		return "mac"
	case "win32", "win", "windows":
		return "windows"
	case "linux", "deb", "appimage":
		return "linux"
	case "linux-rpm", "rpm":
		return "linux-rpm"
	default:
		return p
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
