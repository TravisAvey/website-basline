package routes

import (
	"html/template"
	"net/http"
)

func gallery(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "Gallery Page",
	}

	t, _ := template.ParseFiles("web/templates/pages/gallery.html")
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}
