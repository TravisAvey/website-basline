package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

func parsePostCategories(r *http.Request) ([]database.Category, error) {
	var categories []database.Category

	strCats := strings.Split(r.FormValue("categories"), ",")

	for _, cat := range strCats {
		id, err := database.GetPostCategoryID(cat)
		if err != nil {
			return categories, err
		}
		category := database.Category{
			Category: cat,
			ID:       id,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func parseImageCategories(r *http.Request) ([]database.ImageCategory, error) {
	var categories []database.ImageCategory

	strCats := strings.Split(r.FormValue("categories"), ",")

	for _, cat := range strCats {
		id, err := database.GetGalleryCategoryID(cat)
		if err != nil {
			return categories, err
		}
		category := database.ImageCategory{
			Category: cat,
			ID:       id,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func parsePostForm(r *http.Request) (database.Post, error) {
	err := r.ParseForm()
	if err != nil {
		return database.Post{}, err
	}

	categories, err := parsePostCategories(r)
	if err != nil {
		return database.Post{}, err
	}
	id, err := strconv.ParseInt(r.FormValue("post-id"), 10, 64)
	if err != nil {
		return database.Post{}, err
	}
	post := database.Post{
		Article: database.Article{
			Title:    r.FormValue("title"),
			ID:       id,
			ImageURL: r.FormValue("imageURL"),
			Summary:  r.FormValue("summary"),
			Content:  r.FormValue("content"),
			Slug:     r.FormValue("slug"),
			Keywords: r.FormValue("keywords"),
		},
		Categories: categories,
	}

	return post, nil
}

func parseImageData(r *http.Request) (database.Image, error) {
	err := r.ParseForm()
	if err != nil {
		return database.Image{}, err
	}

	isGallery, err := strconv.ParseBool(r.FormValue("is-gallery"))
	if err != nil {
		return database.Image{}, err
	}

	categories, err := parseImageCategories(r)
	if err != nil {
		return database.Image{}, err
	}

	image := database.Image{
		Image: database.Photo{
			// TODO: get after we upload to S3...
			// ImageURL:  r.FormValue("imageURL"),
			Title:     r.FormValue("title"),
			Summary:   r.FormValue("description"),
			IsGallery: isGallery,
		},
		Categories: categories,
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
	return fmt.Sprintf("%d-%d-%d", y, m, d)
	// return fmt.Sprintf("%d/%d/%d", y, m, d)
}
