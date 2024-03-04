package routes

import (
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, _ *http.Request) {
	// TODO: data from DB for the "home page"
	data := struct {
		Text string
	}{
		Text: "Hello, Templates",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/index.html")

	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
