package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ContactPage renders the contact form, or a thank-you message after submission.
func ContactPage(props PageProps, sent bool) Node {
	props.Title = "Contact — Cold Air Networks"
	props.Description = "Get in touch with Cold Air Networks."

	if sent {
		return page(props,
			Div(Class("mx-auto max-w-xl py-16 text-center"),
				H1(Class("text-3xl font-bold tracking-tight text-slate-900"), Text("Message received")),
				P(Class("mt-4 text-lg text-slate-600"),
					Text("Thanks for reaching out. We'll get back to you within one business day."),
				),
				Div(Class("mt-8"),
					A(Href("/"), Class("text-sm font-medium text-slate-900 underline underline-offset-4"),
						Text("Back to home"),
					),
				),
			),
		)
	}

	return page(props,
		Div(Class("mx-auto max-w-xl py-16"),
			H1(Class("text-3xl font-bold tracking-tight text-slate-900"), Text("Let's talk")),
			P(Class("mt-4 text-lg text-slate-600"),
				Text("Tell us about your project. We'll respond within one business day."),
			),
			Form(Action("/contact"), Method("post"), Class("mt-10 space-y-6"),
				contactField("name", "Name", "text"),
				contactField("email", "Email", "email"),
				Div(
					Label(Class("block text-sm font-medium text-slate-700"), For("message"),
						Text("Message"),
					),
					Textarea(
						Class("mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 shadow-sm focus:border-slate-500 focus:outline-none sm:text-sm"),
						ID("message"), Name("message"), Rows("6"), Required(),
					),
				),
				Button(Type("submit"),
					Class("w-full rounded-lg bg-slate-900 px-6 py-3 text-base font-semibold text-white hover:bg-slate-700"),
					Text("Send message"),
				),
			),
		),
	)
}

func contactField(id, labelText, inputType string) Node {
	return Div(
		Label(Class("block text-sm font-medium text-slate-700"), For(id),
			Text(labelText),
		),
		Input(
			Class("mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 shadow-sm focus:border-slate-500 focus:outline-none sm:text-sm"),
			Type(inputType), ID(id), Name(id), Required(),
		),
	)
}
