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

	t, _ := template.ParseFiles("web/templates/pages/dashboard.html")
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}

// not sure what need endpoints here are needed..
// most CRUD ops already in other routes' go files
//
// Will need to decide how much work here for the baseline
