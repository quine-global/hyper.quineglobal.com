package html

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Release is a single GitHub release with per-platform download URLs.
type Release struct {
	Version     string
	CanaryBuild int    // >0 for canary, e.g. 12 from "v4.0.0-q-canary.12"
	Date        string
	IsCanary    bool
	Assets      map[string]map[string]string // os -> arch -> URL
}

// DownloadURL returns the asset URL for the given OS and arch, or "" if unavailable.
func (r Release) DownloadURL(os, arch string) string {
	if m := r.Assets[os]; m != nil {
		return m[arch]
	}
	return ""
}

// DisplayVersion returns a human-readable version string.
// Canary builds render as "4.0.0 Canary #12".
func (r Release) DisplayVersion() string {
	if r.CanaryBuild > 0 {
		return fmt.Sprintf("%s Canary #%d", r.Version, r.CanaryBuild)
	}
	return r.Version
}

// FallbackReleases is used when the GitHub API is unavailable.
var FallbackReleases = []Release{
	buildFallback("4.2.0", "May 2026", false),
	buildFallback("4.1.0", "Mar 2026", false),
	buildFallback("4.0.0", "Jan 2026", false),
}

func buildFallback(version, date string, canary bool) Release {
	assets := make(map[string]map[string]string)
	for _, os := range []string{"mac", "windows", "linux"} {
		assets[os] = make(map[string]string)
		for _, arch := range []string{"arm64", "x64"} {
			assets[os][arch] = fallbackURL(version, os, arch)
		}
	}
	return Release{Version: version, Date: date, IsCanary: canary, Assets: assets}
}

func fallbackURL(version, os, arch string) string {
	switch os {
	case "windows":
		return fmt.Sprintf("https://github.com/quine-global/hyper/releases/download/v%s/QuineHyper.Setup.%s.exe", version, version)
	case "linux":
		if arch == "arm64" {
			return fmt.Sprintf("https://github.com/quine-global/hyper/releases/download/v%s/QuineHyper-%s-arm64.AppImage", version, version)
		}
		return fmt.Sprintf("https://github.com/quine-global/hyper/releases/download/v%s/QuineHyper-%s.AppImage", version, version)
	default: // mac
		if arch == "arm64" {
			return fmt.Sprintf("https://github.com/quine-global/hyper/releases/download/v%s/QuineHyper-%s-arm64.dmg", version, version)
		}
		return fmt.Sprintf("https://github.com/quine-global/hyper/releases/download/v%s/QuineHyper-%s.dmg", version, version)
	}
}

// DownloadProps carries platform selection and release list for the download page.
type DownloadProps struct {
	DetectedOS   string
	DetectedArch string
	SelectedOS   string
	SelectedArch string
	Releases     []Release
}

type platformOption struct {
	OS    string
	Arch  string
	Label string
}

var platforms = []platformOption{
	{"mac", "arm64", "macOS · Apple Silicon"},
	{"mac", "x64", "macOS · Intel"},
	{"windows", "x64", "Windows · x64"},
	{"linux", "x64", "Linux · x64"},
	{"linux", "arm64", "Linux · ARM64"},
}

func fileExt(os string) string {
	switch os {
	case "windows":
		return ".exe"
	case "linux":
		return ".AppImage"
	case "linux-rpm":
		return ".rpm"
	default:
		return ".dmg"
	}
}

func platformLabel(os, arch string) string {
	effectiveOS := os
	if os == "linux-rpm" {
		effectiveOS = "linux"
	}
	for _, p := range platforms {
		if p.OS == effectiveOS && p.Arch == arch {
			return p.Label
		}
	}
	return effectiveOS + " · " + arch
}

// DownloadPage renders the download page with OS/arch selection and version list.
func DownloadPage(props PageProps, dl DownloadProps) Node {
	props.Title = "Download — Quine Hyper"
	props.Description = "Download Quine Hyper for macOS, Windows, or Linux."

	// Split into stable and canary.
	var stable, canary []Release
	for _, r := range dl.Releases {
		if r.IsCanary {
			canary = append(canary, r)
		} else {
			stable = append(stable, r)
		}
	}

	content := []Node{
		H1(Class("font-mono text-3xl font-bold text-white"), Text("Download Quine Hyper")),
		Div(Class("mt-8"),
			P(Class("font-mono text-xs uppercase tracking-widest text-zinc-500 mb-3"), Text("Platform")),
			Div(Class("flex flex-wrap gap-2"),
				Group(platformButtons(dl)),
			),
			If(dl.SelectedOS == "linux" || dl.SelectedOS == "linux-rpm",
				linuxFormatToggle(dl),
			),
		),
	}
	if len(stable) > 0 {
		content = append(content, stableSection(stable, dl))
	}
	if len(canary) > 0 {
		content = append(content, canarySection(canary, dl))
	}
	content = append(content,
		Div(Class("mt-6 text-center"),
			A(
				Href("https://github.com/quine-global/hyper/releases"),
				Class("font-mono text-xs text-zinc-500 hover:text-white transition-colors underline underline-offset-4"),
				Text("All releases on GitHub ↗"),
			),
		),
	)

	return page(props,
		Div(Class("mx-auto max-w-xl py-16"),
			Group(content),
		),
	)
}

func stableSection(stable []Release, dl DownloadProps) Node {
	latest := stable[0]
	older := stable[1:]
	url := latest.DownloadURL(dl.SelectedOS, dl.SelectedArch)

	return Group{
		// Primary download card
		Div(Class("mt-10 rounded-xl border border-emerald-500/30 bg-zinc-900 p-6"),
			Div(Class("flex items-center justify-between gap-4 flex-wrap"),
				Div(
					Span(Class("font-mono text-xs font-semibold uppercase tracking-widest text-emerald-400"),
						Text("Recommended"),
					),
					H2(Class("mt-1 font-mono text-xl font-bold text-white"),
						Text("Quine Hyper v"+latest.DisplayVersion()),
					),
					P(Class("mt-1 font-mono text-xs text-zinc-500"),
						Text(latest.Date+" · "+platformLabel(dl.SelectedOS, dl.SelectedArch)+fileExt(dl.SelectedOS)),
					),
				),
				downloadButton(url, "bg-emerald-500 hover:bg-emerald-400 text-black"),
			),
		),

		// Older stable releases
		If(len(older) > 0,
			Div(Class("mt-8"),
				P(Class("font-mono text-xs uppercase tracking-widest text-zinc-500 mb-4"), Text("Older releases")),
				Div(Class("divide-y divide-zinc-800 rounded-xl border border-zinc-800"),
					Group(releaseRows(older, dl, "text-emerald-400 hover:text-emerald-300")),
				),
			),
		),
	}
}

func canarySection(canary []Release, dl DownloadProps) Node {
	latest := canary[0]
	older := canary[1:]
	url := latest.DownloadURL(dl.SelectedOS, dl.SelectedArch)

	return Div(Class("mt-10"),
		Div(Class("mb-4 flex items-center gap-3"),
			P(Class("font-mono text-xs uppercase tracking-widest text-yellow-500"), Text("Canary releases")),
			Span(Class("font-mono text-xs text-zinc-600"), Text("— unstable, may contain bugs")),
		),

		// Latest canary card
		Div(Class("rounded-xl border border-yellow-500/30 bg-zinc-900 p-6"),
			Div(Class("flex items-center justify-between gap-4 flex-wrap"),
				Div(
					Span(Class("font-mono text-xs font-semibold uppercase tracking-widest text-yellow-400"),
						Text("Latest Canary"),
					),
					H2(Class("mt-1 font-mono text-xl font-bold text-white"),
						Text("Quine Hyper v"+latest.DisplayVersion()),
					),
					P(Class("mt-1 font-mono text-xs text-zinc-500"),
						Text(latest.Date+" · "+platformLabel(dl.SelectedOS, dl.SelectedArch)+fileExt(dl.SelectedOS)),
					),
				),
				downloadButton(url, "bg-yellow-500 hover:bg-yellow-400 text-black"),
			),
		),

		// Older canary releases
		If(len(older) > 0,
			Div(Class("mt-4 divide-y divide-zinc-800 rounded-xl border border-zinc-800"),
				Group(releaseRows(older, dl, "text-yellow-400 hover:text-yellow-300")),
			),
		),
	)
}

func downloadButton(url, colorCls string) Node {
	if url == "" {
		return Span(
			Class("inline-flex items-center gap-2 rounded-lg border border-zinc-800 px-6 py-3 font-mono text-sm text-zinc-600 shrink-0"),
			Text("Not available"),
		)
	}
	return A(
		Href(url),
		Class("inline-flex items-center gap-2 rounded-lg px-6 py-3 font-mono text-sm font-semibold transition-colors shrink-0 "+colorCls),
		Text("↓ Download"),
	)
}

func releaseRows(releases []Release, dl DownloadProps, linkCls string) []Node {
	nodes := make([]Node, len(releases))
	for i, r := range releases {
		url := r.DownloadURL(dl.SelectedOS, dl.SelectedArch)
		var link Node
		if url != "" {
			link = A(
				Href(url),
				Class("font-mono text-xs transition-colors "+linkCls),
				Text("↓ "+platformLabel(dl.SelectedOS, dl.SelectedArch)+fileExt(dl.SelectedOS)),
			)
		} else {
			link = Span(
				Class("font-mono text-xs text-zinc-600"),
				Text("Not available"),
			)
		}
		nodes[i] = Div(Class("flex items-center justify-between px-5 py-4"),
			Div(Class("flex items-center gap-4"),
				Span(Class("font-mono text-sm font-semibold text-white"), Text("v"+r.DisplayVersion())),
				Span(Class("font-mono text-xs text-zinc-500"), Text(r.Date)),
			),
			link,
		)
	}
	return nodes
}

func platformButtons(dl DownloadProps) []Node {
	nodes := make([]Node, len(platforms))
	// Treat linux-rpm as linux for platform button highlight/detection purposes.
	effectiveSelectedOS := dl.SelectedOS
	if dl.SelectedOS == "linux-rpm" {
		effectiveSelectedOS = "linux"
	}
	for i, p := range platforms {
		isSelected := p.OS == effectiveSelectedOS && p.Arch == dl.SelectedArch
		isDetected := p.OS == dl.DetectedOS && p.Arch == dl.DetectedArch

		// Preserve the linux format (appimage vs rpm) when switching arch.
		targetOS := p.OS
		if p.OS == "linux" && dl.SelectedOS == "linux-rpm" {
			targetOS = "linux-rpm"
		}

		var cls string
		if isSelected {
			cls = "inline-flex items-center gap-1.5 rounded-md bg-emerald-500 px-3 py-1.5 font-mono text-xs font-semibold text-black"
		} else {
			cls = "inline-flex items-center gap-1.5 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-1.5 font-mono text-xs text-zinc-300 hover:border-zinc-500 hover:text-white transition-colors"
		}

		children := []Node{Text(p.Label)}
		if isDetected && !isSelected {
			children = append(children,
				Span(Class("rounded bg-zinc-700 px-1 py-0.5 text-zinc-400 text-[10px]"), Text("auto")),
			)
		} else if isDetected && isSelected {
			children = append(children,
				Span(Class("rounded bg-black/20 px-1 py-0.5 text-[10px]"), Text("auto")),
			)
		}

		nodes[i] = A(
			Href(fmt.Sprintf("/download?os=%s&arch=%s", targetOS, p.Arch)),
			Class(cls),
			Group(children),
		)
	}
	return nodes
}

func linuxFormatToggle(dl DownloadProps) Node {
	isRPM := dl.SelectedOS == "linux-rpm"
	appImageHref := fmt.Sprintf("/download?os=linux&arch=%s", dl.SelectedArch)
	rpmHref := fmt.Sprintf("/download?os=linux-rpm&arch=%s", dl.SelectedArch)

	activeCls := "font-mono text-xs font-semibold text-white"
	inactiveCls := "font-mono text-xs text-zinc-500 hover:text-zinc-300 transition-colors"

	appImageCls := activeCls
	rpmCls := inactiveCls
	if isRPM {
		appImageCls = inactiveCls
		rpmCls = activeCls
	}

	return Div(Class("mt-4"),
		P(Class("font-mono text-xs uppercase tracking-widest text-zinc-500 mb-3"), Text("Format")),
		Div(Class("flex items-center gap-2"),
			A(Href(appImageHref), Class(appImageCls), Text("AppImage")),
			Span(Class("font-mono text-xs text-zinc-700"), Text("|")),
			A(Href(rpmHref), Class(rpmCls), Text("RPM")),
		),
	)
}
