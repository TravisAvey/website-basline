package routes

import (
	"html/template"
	"net/http"
)

func about(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "About Page",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/about.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// not much to do here.. maybe just update the template with all that
// is needed on the about page...
