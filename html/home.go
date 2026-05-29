package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomePage is the Quine Hyper landing page.
func HomePage(props PageProps) Node {
	props.Title = "Quine Hyper — The thoughtful terminal"
	props.Description = "Quine Hyper is a beautiful, extensible terminal emulator forked from Vercel's Hyper. Built on web technology with a thriving plugin ecosystem."

	return page(props,
		// Hero
		Div(Class("py-20 text-center"),
			Div(Class("mb-6 inline-flex items-center gap-2 rounded-full border border-zinc-700 bg-zinc-900 px-4 py-1.5"),
				Span(Class("h-2 w-2 rounded-full bg-violet-400")),
				Span(Class("font-mono text-xs text-zinc-400"), Text("Forked from Vercel's Hyper")),
			),
			H1(Class("font-mono text-5xl font-bold tracking-tight text-white sm:text-6xl"),
				Text("The thoughtful terminal."),
			),
			P(Class("mx-auto mt-6 max-w-2xl font-mono text-lg text-zinc-400"),
				Text("A beautiful, extensible terminal built with care."),
			),
			Div(Class("mt-10 flex flex-wrap items-center justify-center gap-4"),
				A(Href("/download"),
					Class("inline-flex items-center gap-2 rounded-lg bg-violet-500 px-8 py-3 font-mono text-base font-semibold text-black hover:bg-violet-400 transition-colors"),
					Text("Download for free"),
				),
				A(Href("https://github.com/quine-global/hyper"),
					Class("inline-flex items-center gap-2 rounded-lg border border-zinc-700 bg-zinc-900 px-8 py-3 font-mono text-base font-semibold text-white hover:border-zinc-500 transition-colors"),
					Text("View on GitHub"),
				),
			),
		),

		// Terminal window mockup
		Div(Class("mx-auto mb-20 max-w-3xl overflow-hidden rounded-xl border border-zinc-800 bg-zinc-900 shadow-2xl"),
			// Title bar
			Div(Class("flex items-center gap-2 border-b border-zinc-800 bg-zinc-900 px-4 py-3"),
				Div(Class("h-3 w-3 rounded-full bg-red-500 opacity-80")),
				Div(Class("h-3 w-3 rounded-full bg-yellow-500 opacity-80")),
				Div(Class("h-3 w-3 rounded-full bg-violet-500 opacity-80")),
				Span(Class("ml-4 font-mono text-xs text-zinc-500"), Text("quine-hyper — bash — 80×24")),
			),
			// Terminal body
			Div(Class("p-6 font-mono text-sm leading-relaxed"),
				termPrompt("deep-thought --answer"),
				termOutput("42"),
				termPrompt("hyper install infinite-improbability"),
				termOutput("✔ Engaging drive..."),
				termOutput("✔ Probability: 2:1 against."),
				termOutput("✔ Don't Panic."),
				termPrompt("ls ~/projects"),
				termOutput("drwxr-xr-x  earth/"),
				termOutput("drwxr-xr-x  magrathea/"),
				termOutput("drwxr-xr-x  slartibartfast/"),
				Div(Class("flex items-center gap-0 text-violet-400"),
					Span(Class("text-zinc-500"), Text("❯ ")),
					Span(Class("animate-pulse"), Text("█")),
				),
			),
		),

		// Features
		Div(Class("pb-20"),
			H2(Class("text-center font-mono text-xs font-semibold uppercase tracking-widest text-zinc-500"),
				Text("Why Quine Hyper"),
			),
			Div(Class("mt-12 grid grid-cols-1 gap-6 md:grid-cols-2"),
				featureCard("Plugin Ecosystem",
					Group{
						Text("Install plugins and themes from npm. Browse "),
						A(Href("https://hyper.is/plugins"),
							Class("text-violet-400 hover:text-violet-300 underline"),
							Text("featured plugins"),
						),
						Text(" and "),
						A(Href("https://hyper.is/themes"),
							Class("text-violet-400 hover:text-violet-300 underline"),
							Text("featured themes"),
						),
						Text(". Build your own with HTML, CSS, and JavaScript."),
					},
					"npm install hyper-snazzy",
				),
				featureCard("Cross-Platform",
					Text("Runs on macOS, Windows, and Linux. One consistent experience across every machine you use."),
					"Works everywhere you do.",
				),
				featureCard("Web-Native",
					Text("Built on Electron, React, and xterm.js. Hack on the terminal itself using the same tools you use every day."),
					"HTML + CSS + JS = your terminal",
				),
				featureCard("Fully Customizable",
					Text("Tweak every detail with CSS. Change fonts, colors, animations, and layouts — no config file limits."),
					"hyper.json → your rules",
				),
			),
		),

		// Origin story
		Div(Class("mb-20 rounded-xl border border-zinc-800 bg-zinc-900 p-10 text-center"),
			Img(Src("/images/quineglobal-logo.png"), Alt("Quine Global"), Class("mx-auto mb-6 h-12 w-12 rounded-full")),
			H2(Class("font-mono text-2xl font-bold text-white"), Text("Built on Hyper")),
			P(Class("mx-auto mt-4 max-w-xl font-mono text-base text-zinc-400"),
				Text("Quine Hyper is a fork of "),
				A(Href("https://hyper.is"), Class("text-violet-400 hover:text-violet-300 underline"), Text("Hyper")),
				Text(", the open-source terminal by Vercel. We're building on that foundation with new capabilities and active development."),
			),
			Div(Class("mt-8"),
				A(Href("https://github.com/quine-global/hyper"),
					Class("inline-block rounded-lg border border-zinc-700 bg-zinc-800 px-6 py-3 font-mono text-sm font-semibold text-white hover:border-zinc-500 transition-colors"),
					Text("Star us on GitHub"),
				),
			),
		),
	)
}

func termPrompt(cmd string) Node {
	return Div(Class("flex items-center gap-2"),
		Span(Class("text-violet-400"), Text("❯")),
		Span(Class("text-white"), Text(cmd)),
	)
}

func termOutput(line string) Node {
	return Div(Class("pl-5 text-zinc-400"), Text(line))
}

func featureCard(title string, body Node, code string) Node {
	return Div(Class("rounded-xl border border-zinc-800 bg-zinc-900 p-8"),
		H3(Class("font-mono text-base font-semibold text-white"), Text(title)),
		P(Class("mt-3 font-mono text-sm text-zinc-400 leading-relaxed"), body),
		Div(Class("mt-5 inline-block rounded-md bg-zinc-800 px-3 py-1.5"),
			Span(Class("font-mono text-xs text-violet-400"), Text(code)),
		),
	)
}
