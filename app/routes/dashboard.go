package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	page := views.Page("Dashboard Page")
	page.Render(r.Context(), w)
}
