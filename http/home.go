package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrz1836/postmark"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"app/html"
)

const contactRecipient = "info@coldairnetworks.com"

// Home registers the index and contact handlers.
func (s *Server) Home(r chi.Router) {
	r.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.HomePage(html.PageProps{}), nil
	}))

	r.Get("/contact", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		sent := r.URL.Query().Get("sent") == "1"
		return html.ContactPage(html.PageProps{}, sent), nil
	}))

	r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		_, err := s.postmark.SendEmail(r.Context(), postmark.Email{
			From:     contactRecipient,
			To:       contactRecipient,
			ReplyTo:  email,
			Subject:  fmt.Sprintf("Contact form: %s", name),
			TextBody: fmt.Sprintf("Name: %s\nEmail: %s\n\n%s", name, email, message),
		})
		if err != nil {
			s.log.Error("Failed to send contact email", "error", err)
			http.Error(w, "Failed to send message, please try again later.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/contact?sent=1", http.StatusSeeOther)
	})
}
