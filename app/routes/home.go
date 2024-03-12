package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/travisavey/baseline/app/logging"
)

func index(w http.ResponseWriter, r *http.Request) {
	logging.LogAccess(fmt.Sprintf("index page requested by: %s", r.RemoteAddr))
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
