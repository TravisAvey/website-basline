package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"html/template"
	"net/http"
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Dashboard Page",
	}

	files := []string{"web/templates/base.html", "web/templates/pages/dashboard.html"}
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// not sure what need endpoints here are needed..
// most CRUD ops already in other routes' go files
//
// Will need to decide how much work here for the baseline
