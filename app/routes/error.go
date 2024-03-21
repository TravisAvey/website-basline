package routes

import (
	"html/template"
	"net/http"
)

type errMsg struct {
	Message   string
	Title     string
	ImageURL  string
	ErrorCode int
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	msg := errMsg{
		ErrorCode: 404,
		Message:   "Sorry, the page requested was not found",
		Title:     "_Not Found",
		ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
	}

	sendErrorTemplate(msg, w)
}

func notAllowed(w http.ResponseWriter, _ *http.Request) {
	msg := errMsg{
		ErrorCode: 404,
		Message:   "Sorry, the method requested is not allowed",
		Title:     "_Not Allowed",
		ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
	}
	sendErrorTemplate(msg, w)
}

func sendErrorTemplate(msg errMsg, w http.ResponseWriter) {
	files := getBaseTemplates()
	files = append(files, "web/templates/pages/error/error.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", msg)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
