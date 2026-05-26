// Package HTML holds all the common HTML components and utilities.
package html

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

var hashOnce sync.Once
var appCSSPath string

// PageProps are properties for the [page] component.
type PageProps struct {
	Title       string
	Description string
}

// page layout with nav, main content, and footer.
func page(props PageProps, children ...Node) Node {
	hashOnce.Do(func() {
		appCSSPath = getHashedPath("public/styles/app.css")
	})

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href(appCSSPath)),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		},
		Body: Group{
			siteNav(),
			Main(Class("mx-auto max-w-5xl px-6 py-12"),
				Group(children),
			),
			siteFooter(),
		},
	})
}

func siteNav() Node {
	return Header(Class("border-b border-slate-100"),
		Div(Class("mx-auto flex max-w-5xl items-center justify-between px-6 py-4"),
			A(Href("/"), Class("text-lg font-bold tracking-tight text-slate-900"),
				Text("Cold Air Networks"),
			),
			A(Href("/contact"), Class("rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-700"),
				Text("Contact us"),
			),
		),
	)
}

func siteFooter() Node {
	return Footer(Class("mt-24 border-t border-slate-100 py-8"),
		Div(Class("mx-auto flex max-w-5xl items-center justify-between px-6 text-sm text-slate-500"),
			Text("© 2025 Cold Air Networks"),
			Div(Class("flex gap-6"),
				A(Href("/"), Class("hover:text-slate-900"), Text("Home")),
				A(Href("/contact"), Class("hover:text-slate-900"), Text("Contact")),
			),
		),
	)
}

func getHashedPath(path string) string {
	externalPath := strings.TrimPrefix(path, "public")
	ext := filepath.Ext(path)
	if ext == "" {
		panic("no extension found")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("%v.x%v", strings.TrimSuffix(externalPath, ext), ext)
	}

	return fmt.Sprintf("%v.%x%v", strings.TrimSuffix(externalPath, ext), sha256.Sum256(data), ext)
}
