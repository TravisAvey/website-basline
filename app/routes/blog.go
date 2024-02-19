package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/views"
)

func blog(w http.ResponseWriter, r *http.Request) {
	page := views.Page("Blog Page")
	page.Render(r.Context(), w)
}
