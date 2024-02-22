package routes

import (
	"net/http"

	"github.com/travisavey/baseline/app/database"
)

func getBaseTemplates() []string {
	return []string{
		"web/templates/base.html",
		"web/templates/partials/header.html",
		"web/templates/partials/nav.html",
		"web/templates/partials/footer.html",
	}
}

func parseFormData(r *http.Request) (database.Post, error) {
	err := r.ParseForm()
	if err != nil {
		return database.Post{}, err
	}

	post := database.Post{
		Article: database.Article{
			Title:    r.FormValue("title"),
			ImageURL: r.FormValue("imageURL"),
			Summary:  r.FormValue("summary"),
			Content:  r.FormValue("content"),
			Slug:     r.FormValue("slug"),
			Keywords: r.FormValue("keywords"),
		},
		// TODO: figure this out..
		// should have our form have pre configured
		// selectable categories...
		// also, figure out how to parse an array here
		Categories: []database.Category{
			{
				ID:       1,
				Category: r.FormValue("category"),
			},
		},
	}

	return post, nil
}
