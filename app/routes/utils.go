package routes

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/travisavey/baseline/app/database"
)

type ResponseType int

const (
	Info ResponseType = iota
	Warn
	Success
	Error
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

	catID, err := strconv.ParseInt(r.FormValue("categoryId"), 10, 64)
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
				ID: catID,
			},
		},
	}

	return post, nil
}

func parseImageData(r *http.Request) (database.Image, error) {
	err := r.ParseForm()
	if err != nil {
		return database.Image{}, err
	}

	isGallery, err := strconv.ParseBool(r.FormValue("isGallery"))
	if err != nil {
		return database.Image{}, err
	}

	catID, err := strconv.ParseUint(r.FormValue("categoryId"), 10, 64)
	if err != nil {
		return database.Image{}, err
	}

	image := database.Image{
		Image: database.Photo{
			ImageURL:  r.FormValue("imageURL"),
			Title:     r.FormValue("title"),
			Summary:   r.FormValue("summary"),
			IsGallery: isGallery,
		},
		Categories: []database.ImageCategory{
			{
				ID: catID,
			},
		},
	}

	return image, nil
}

func sendResponseMsg(msg string, res ResponseType, w http.ResponseWriter) error {
	var resType string
	switch res {
	case Warn:
		resType = "warn"
	case Info:
		resType = "info"
	case Success:
		resType = "success"
	case Error:
		resType = "error"
	}
	data := struct {
		Message string
		Type    string
	}{
		Message: msg,
		Type:    resType,
	}

	t, _ := template.ParseFiles("web/templates/responses/message.html")
	return t.Execute(w, data)
}
