package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"html/template"
	"net/http"
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "Dashboard Page",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/index.html")

	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}
