package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func index(w http.ResponseWriter, r *http.Request) {
	page := views.Page("Hello, index!")
	page.Render(r.Context(), w)
}
