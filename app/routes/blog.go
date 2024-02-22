package routes

import (
	"html/template"
	"net/http"
)

func blog(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Text string
	}{
		Text: "Blog Page",
	}

	t, _ := template.ParseFiles("web/templates/pages/blog.html")
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte("Error processing templates.."))
	}
}

// create a post

// get a single post

// get all posts

// update a post

// delete a post
