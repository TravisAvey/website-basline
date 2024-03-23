package routes

import (
	"html/template"
	"net/http"
)

func termsOfUse(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text     string
		ImageURL string
	}{
		Text:     "Terms of Use",
		ImageURL: "https://picsum.photos/1920/1080/?blur=2",
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/legal/terms.html")
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
