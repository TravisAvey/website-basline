package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func contact(w http.ResponseWriter, r *http.Request) {
	page := views.Page("Contact Page")
	page.Render(r.Context(), w)
}
