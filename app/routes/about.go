package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func about(w http.ResponseWriter, r *http.Request) {
	page := views.Page("About Page")
	page.Render(r.Context(), w)
}
