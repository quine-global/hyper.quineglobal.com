package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"app/html"
)

// Home registers the index and download handlers.
func (s *Server) Home(r chi.Router) {
	r.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.HomePage(html.PageProps{}), nil
	}))

	r.Get("/download", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// On-demand refresh, debounced to once per 2 minutes across all requests.
		go s.releases.TryRefresh(context.Background(), 2*time.Minute)

		detectedOS, detectedArch := detectPlatform(r.Header.Get("User-Agent"))

		selectedOS := detectedOS
		selectedArch := detectedArch

		if p := r.URL.Query().Get("os"); p == "mac" || p == "windows" || p == "linux" || p == "linux-rpm" || p == "linux-deb" || p == "linux-snap" || p == "linux-pacman" {
			selectedOS = p
		}
		if a := r.URL.Query().Get("arch"); a == "arm64" || a == "x64" {
			selectedArch = a
		}

		return html.DownloadPage(html.PageProps{}, html.DownloadProps{
			DetectedOS:   detectedOS,
			DetectedArch: detectedArch,
			SelectedOS:   selectedOS,
			SelectedArch: selectedArch,
			Releases:     s.releases.Get(),
		}), nil
	}))
}

func detectPlatform(ua string) (os, arch string) {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "windows"):
		return "windows", "x64"
	case strings.Contains(u, "macintosh"), strings.Contains(u, "mac os"):
		return "mac", "arm64" // default Apple Silicon; user can override
	case strings.Contains(u, "linux"):
		return "linux", "x64"
	default:
		return "mac", "arm64"
	}
}
