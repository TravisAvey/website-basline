package routes

import (
	"html/template"
	"net/http"
)

func contact(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "Contact Page",
	}

	t, _ := template.ParseFiles("web/templates/pages/contact.html")
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
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
