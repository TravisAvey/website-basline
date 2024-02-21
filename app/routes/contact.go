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

	files := getBaseTemplates()
	files = append(files, "web/templates/index.html")

	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}
