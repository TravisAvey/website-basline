package routes

import (
	"html/template"
	"net/http"
)

func notFound(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Message   string
		Title     string
		ImageURL  string
		ErrorCode int
	}{
		ErrorCode: 404,
		Message:   "Sorry, the page requested was not found",
		Title:     "_Not Found",
		ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/error/error.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
