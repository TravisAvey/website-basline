package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"fmt"
	"net/http"
	"text/template"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Text string
	}{
		Text: "<h1>Dashboard</h1>",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/index.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		// TODO: log error
		fmt.Println("Error executing template: ", err.Error())
	}
}
