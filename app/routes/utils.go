package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

func getResponseMsg(msg string, res ResponseType) string {
	var buf bytes.Buffer
	data := struct {
		Message string
	}{
		Message: msg,
	}

	var t *template.Template
	var err error

	switch res {
	case Warn:
		t, err = template.ParseFiles("web/templates/messages/warn.html")
	case Info:
		t, err = template.ParseFiles("web/templates/messages/info.html")
	case Success:
		t, err = template.ParseFiles("web/templates/messages/success.html")
	case Error:
		t, err = template.ParseFiles("web/templates/messages/error.html")
	}

	if err != nil {
		return ""
	}

	err = t.Execute(&buf, data)
	if err != nil {
		return ""
	}

	return buf.String()
}

func sendResponseMsg(msg string, res ResponseType, w http.ResponseWriter) error {
	data := struct {
		Message string
	}{
		Message: msg,
	}

	var t *template.Template
	var err error

	switch res {
	case Warn:
		t, err = template.ParseFiles("web/templates/messages/warn.html")
	case Info:
		t, err = template.ParseFiles("web/templates/messages/info.html")
	case Success:
		t, err = template.ParseFiles("web/templates/messages/success.html")
	case Error:
		t, err = template.ParseFiles("web/templates/messages/error.html")
	}

	if err != nil {
		return err
	} else {
		return t.Execute(w, data)
	}
}

func parseDate(time time.Time) string {
	y, m, d := time.Date()
	return fmt.Sprintf("%d/%d/%d", y, m, d)
}
