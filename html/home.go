package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomePage is the consultancy landing page.
func HomePage(props PageProps) Node {
	props.Title = "Cold Air Networks — Web Consulting"
	props.Description = "We build reliable web systems for businesses that depend on them."

	return page(props,
		Div(Class("py-24 text-center"),
			H1(Class("text-5xl font-bold tracking-tight text-slate-900 sm:text-6xl"),
				Text("Software that holds up."),
			),
			P(Class("mx-auto mt-6 max-w-2xl text-xl text-slate-600"),
				Text("Cold Air Networks builds web systems for businesses that need them to work. Architecture, development, and delivery — without the noise."),
			),
			Div(Class("mt-10"),
				A(Href("/contact"),
					Class("inline-block rounded-lg bg-slate-900 px-8 py-3 text-base font-semibold text-white hover:bg-slate-700"),
					Text("Start a conversation"),
				),
			),
		),
		Hr(Class("border-slate-200")),
		Div(Class("py-20"),
			H2(Class("text-center text-sm font-semibold uppercase tracking-widest text-slate-500"),
				Text("What we do"),
			),
			Div(Class("mt-12 grid grid-cols-1 gap-8 md:grid-cols-3"),
				serviceCard("Web Development", "Full-stack applications built to last. Clean code, sound architecture, no technical debt handed off to you."),
				serviceCard("System Design", "We design systems that fit your scale today and grow with you. No over-engineering, no shortcuts."),
				serviceCard("Technical Consulting", "Straight advice on hard decisions. We help you pick the right tools and avoid expensive mistakes."),
			),
		),
	)
}

func serviceCard(title, body string) Node {
	return Div(Class("rounded-xl border border-slate-200 p-8"),
		H3(Class("text-lg font-semibold text-slate-900"), Text(title)),
		P(Class("mt-3 text-base text-slate-600"), Text(body)),
	)
}
