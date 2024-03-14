package routes

import (
	"html/template"
	"net/http"
)

func contact(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text     string
		ImageURL string
	}{
		Text:     "Contact Page",
		ImageURL: "https://picsum.photos/1920/1080/?blur=2",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/contact.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// endpoint for a post request for a contact here
func contactForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	sendResponseMsg("Message sent successfully.", Success, w)
}

//
// TODO: think how this will work: fire off an email
// what email service to setup and use in the backend here?
// sendgrid, twilio, mailchimp, etc
