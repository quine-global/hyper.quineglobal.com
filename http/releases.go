package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"app/html"
)

const releasesAPI = "https://api.github.com/repos/quine-global/hyper/releases?per_page=30"

// ReleaseCache is a thread-safe cache of GitHub releases.
type ReleaseCache struct {
	log         *slog.Logger
	mu          sync.RWMutex
	releases    []html.Release
	lastFetched time.Time
	fetchMu     sync.Mutex // TryLock prevents concurrent fetches
}

// Get returns cached releases, falling back to hardcoded data if cache is empty.
func (c *ReleaseCache) Get() []html.Release {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.releases) == 0 {
		return []html.Release{}
	}
	return c.releases
}

// ForceRefresh does a blocking, unconditional fetch regardless of cache age.
// It waits for any in-progress fetch to finish before starting its own.
func (c *ReleaseCache) ForceRefresh(ctx context.Context) error {
	c.fetchMu.Lock()
	defer c.fetchMu.Unlock()

	releases, err := fetchReleases(ctx)
	if err != nil {
		c.log.Warn("Force refresh failed", "error", err)
		return err
	}
	if len(releases) == 0 {
		return nil
	}

	c.mu.Lock()
	c.releases = releases
	c.lastFetched = time.Now()
	c.mu.Unlock()

	c.log.Info("Force-refreshed releases from GitHub", "count", len(releases))
	return nil
}

// TryRefresh fetches fresh releases if the cache is older than minAge.
// If a fetch is already in progress it returns immediately (non-blocking).
func (c *ReleaseCache) TryRefresh(ctx context.Context, minAge time.Duration) {
	if !c.fetchMu.TryLock() {
		return
	}
	defer c.fetchMu.Unlock()

	c.mu.RLock()
	stale := time.Since(c.lastFetched) >= minAge
	c.mu.RUnlock()

	if !stale {
		return
	}

	releases, err := fetchReleases(ctx)
	if err != nil {
		c.log.Warn("Failed to fetch releases from GitHub", "error", err)
		return
	}
	if len(releases) == 0 {
		return // keep existing cache rather than replacing with empty
	}

	c.mu.Lock()
	c.releases = releases
	c.lastFetched = time.Now()
	c.mu.Unlock()

	c.log.Info("Refreshed releases from GitHub", "count", len(releases))
}

type ghRelease struct {
	TagName     string    `json:"tag_name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func fetchReleases(ctx context.Context) ([]html.Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	releases := make([]html.Release, 0, len(raw))
	for _, gr := range raw {
		if gr.Draft {
			continue
		}
		version, canaryBuild := parseTag(strings.TrimPrefix(gr.TagName, "v"))
		r := html.Release{
			TagName:     gr.TagName,
			Version:     version,
			CanaryBuild: canaryBuild,
			Date:        gr.PublishedAt.Format("Jan 2006"),
			PublishedAt: gr.PublishedAt,
			IsCanary:    gr.Prerelease || canaryBuild > 0,
			Assets:      make(map[string]map[string]string),
		}
		for _, a := range gr.Assets {
			os, arch, ok := classifyAsset(a.Name)
			if !ok {
				continue
			}
			if r.Assets[os] == nil {
				r.Assets[os] = make(map[string]string)
			}
			// Keep first match per os/arch (most specific asset listed first by GitHub)
			if _, exists := r.Assets[os][arch]; !exists {
				r.Assets[os][arch] = a.BrowserDownloadURL
			}
		}
		releases = append(releases, r)
	}
	return releases, nil
}

// parseTag splits a tag like "4.0.0-q-canary.12" into base version "4.0.0" and build number 12.
// Returns build 0 for stable tags.
func parseTag(tag string) (version string, canaryBuild int) {
	// Canary: "4.0.0-q-canary.12" → ("4.0.0", 12)
	const canaryMarker = "-q-canary."
	if idx := strings.Index(tag, canaryMarker); idx != -1 {
		if n, err := strconv.Atoi(tag[idx+len(canaryMarker):]); err == nil {
			return tag[:idx], n
		}
	}
	// Stable Quine builds: "4.0.0-q" or "4.0.0q" → "4.0.0"
	tag = strings.TrimSuffix(tag, "-q")
	tag = strings.TrimSuffix(tag, "q")
	return tag, 0
}

// classifyAsset maps an asset filename to an OS and architecture.
func classifyAsset(name string) (os, arch string, ok bool) {
	l := strings.ToLower(name)
	switch {
	case strings.HasSuffix(l, ".exe"):
		return "windows", "x64", true
	case strings.Contains(l, "arm64") && strings.HasSuffix(l, ".dmg"):
		return "mac", "arm64", true
	case strings.HasSuffix(l, ".dmg"):
		return "mac", "x64", true
	case strings.Contains(l, "arm64") && strings.HasSuffix(l, ".appimage"):
		return "linux", "arm64", true
	case strings.HasSuffix(l, ".appimage"):
		return "linux", "x64", true
	case strings.Contains(l, "arm64") && strings.HasSuffix(l, ".rpm"):
		return "linux-rpm", "arm64", true
	case strings.HasSuffix(l, ".rpm"):
		return "linux-rpm", "x64", true
	default:
		return "", "", false
	}
}
