package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func gallery(w http.ResponseWriter, r *http.Request) {
	page := views.Page("Gallery Page")
	page.Render(r.Context(), w)
}
