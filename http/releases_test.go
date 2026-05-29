package http

import (
	"testing"

	"maragu.dev/is"
)

func TestClassifyAsset(t *testing.T) {
	cases := []struct {
		name       string
		wantOS     string
		wantArch   string
		wantOK     bool
	}{
		// Windows
		{"QuineHyper.Setup.4.0.0.exe", "windows", "x64", true},
		{"QuineHyper.Setup.4.0.0-q.exe", "windows", "x64", true},

		// macOS
		{"QuineHyper-4.0.0.dmg", "mac", "x64", true},
		{"QuineHyper-4.0.0-arm64.dmg", "mac", "arm64", true},

		// Linux AppImage x64
		{"QuineHyper-4.0.0.AppImage", "linux", "x64", true},
		{"quinehyper-4.0.0.appimage", "linux", "x64", true},

		// Linux AppImage arm64
		{"QuineHyper-4.0.0-arm64.AppImage", "linux", "arm64", true},

		// Linux AppImage armv7l — must NOT be classified as x64
		{"QuineHyper-4.0.0-armv7l.AppImage", "", "", false},
		{"quinehyper-4.0.0-armv7l.appimage", "", "", false},

		// Linux RPM x64 — real asset uses x86_64 suffix
		{"QuineHyper-4.0.0.rpm", "linux-rpm", "x64", true},
		{"Hyper-4.0.0-x86_64.rpm", "linux-rpm", "x64", true},

		// Linux RPM arm64 — real asset uses aarch64, not arm64
		{"QuineHyper-4.0.0-arm64.rpm", "linux-rpm", "arm64", true},
		{"Hyper-4.0.0-aarch64.rpm", "linux-rpm", "arm64", true},

		// Linux RPM armv7l — must NOT be classified as x64
		{"Hyper-4.0.0-armv7l.rpm", "", "", false},

		// Linux DEB x64 — real asset uses amd64 suffix
		{"Hyper-4.0.0-q-canary.8-amd64.deb", "linux-deb", "x64", true},

		// Linux DEB arm64
		{"Hyper-4.0.0-q-canary.8-arm64.deb", "linux-deb", "arm64", true},

		// Linux DEB armv7l — must NOT be classified as x64
		{"Hyper-4.0.0-q-canary.8-armv7l.deb", "", "", false},

		// Linux Snap — x64 only
		{"Hyper-4.0.0-q-canary.8-amd64.snap", "linux-snap", "x64", true},

		// Linux Pacman x64
		{"Hyper-4.0.0-q-canary.8-x64.pacman", "linux-pacman", "x64", true},

		// Linux Pacman arm64 — real asset uses aarch64 suffix
		{"Hyper-4.0.0-q-canary.8-aarch64.pacman", "linux-pacman", "arm64", true},

		// Linux Pacman armv7l — must NOT be classified as x64
		{"Hyper-4.0.0-q-canary.8-armv7l.pacman", "", "", false},

		// Unknown / unrecognised
		{"QuineHyper-4.0.0.zip", "", "", false},
		{"latest.yml", "", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotOS, gotArch, gotOK := classifyAsset(c.name)
			is.Equal(t, c.wantOK, gotOK)
			is.Equal(t, c.wantOS, gotOS)
			is.Equal(t, c.wantArch, gotArch)
		})
	}
}

func TestParseTag(t *testing.T) {
	cases := []struct {
		tag        string
		wantVer    string
		wantBuild  int
	}{
		{"4.0.0", "4.0.0", 0},
		{"4.0.0-q", "4.0.0", 0},
		{"4.0.0q", "4.0.0", 0},
		{"4.0.0-q-canary.12", "4.0.0", 12},
		{"4.0.0-q-canary.1", "4.0.0", 1},
	}

	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			gotVer, gotBuild := parseTag(c.tag)
			is.Equal(t, c.wantVer, gotVer)
			is.Equal(t, c.wantBuild, gotBuild)
		})
	}
}
