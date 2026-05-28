package http

import (
	"testing"

	"maragu.dev/is"

	"app/html"
)

func TestSemverCmp(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"4.0.0", "4.0.0", 0},
		{"4.1.0", "4.0.0", 1},
		{"4.0.0", "4.1.0", -1},
		{"4.0.1", "4.0.0", 1},
		{"5.0.0", "4.9.9", 1},
		{"4.10.0", "4.9.0", 1},
	}
	for _, c := range cases {
		t.Run(c.a+"_vs_"+c.b, func(t *testing.T) {
			is.Equal(t, c.want, semverCmp(c.a, c.b))
		})
	}
}

func TestNewerCanary(t *testing.T) {
	rel := func(base string, build int) html.Release {
		return html.Release{Version: base, CanaryBuild: build, IsCanary: true}
	}

	t.Run("higher base is newer for canary client", func(t *testing.T) {
		is.True(t, newerCanary(rel("4.1.0", 1), "4.0.0", 5, true))
	})
	t.Run("same base higher build is newer for canary client", func(t *testing.T) {
		is.True(t, newerCanary(rel("4.0.0", 6), "4.0.0", 5, true))
	})
	t.Run("same base same build is not newer for canary client", func(t *testing.T) {
		is.True(t, !newerCanary(rel("4.0.0", 5), "4.0.0", 5, true))
	})
	t.Run("lower base is not newer for canary client", func(t *testing.T) {
		is.True(t, !newerCanary(rel("3.9.0", 99), "4.0.0", 1, true))
	})
	t.Run("same base canary is newer than stable client", func(t *testing.T) {
		is.True(t, newerCanary(rel("4.0.0", 1), "4.0.0", 0, false))
	})
	t.Run("higher base canary is newer than stable client", func(t *testing.T) {
		is.True(t, newerCanary(rel("4.1.0", 1), "4.0.0", 0, false))
	})
	t.Run("lower base canary is not newer than stable client", func(t *testing.T) {
		is.True(t, !newerCanary(rel("3.9.0", 99), "4.0.0", 0, false))
	})
}

func TestNormalizePlatformForAssets(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"darwin", "mac"},
		{"Darwin", "mac"},
		{"mac", "mac"},
		{"win32", "windows"},
		{"windows", "windows"},
		{"deb", "linux"},
		{"linux", "linux"},
		{"linux-rpm", "linux-rpm"},
		{"rpm", "linux-rpm"},
		{"", ""},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			is.Equal(t, c.want, normalizePlatformForAssets(c.in))
		})
	}
}

func TestFirstNonEmpty(t *testing.T) {
	is.Equal(t, "a", firstNonEmpty("a", "b"))
	is.Equal(t, "b", firstNonEmpty("", "b"))
	is.Equal(t, "", firstNonEmpty("", ""))
}
