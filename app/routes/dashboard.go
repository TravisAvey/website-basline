package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"html/template"
	"net/http"

	"github.com/travisavey/baseline/app/database"
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Dashboard Page",
	}

	files := []string{"web/templates/base.html", "web/templates/pages/dashboard.html"}
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// not sure what need endpoints here are needed..
// most CRUD ops already in other routes' go files
//
// Will need to decide how much work here for the baseline

func dashboardPosts(w http.ResponseWriter, _ *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		sendResponseMsg("Failed to get blog posts", Error, w)
		return
	}

	for i := range posts {
		posts[i].Article.PostedStr = parseDate(posts[i].Article.DatePosted.Time)
		posts[i].Article.UpdatedStr = parseDate(posts[i].Article.DateUpdated.Time)
	}
	numPosts := len(posts)

	data := struct {
		Posts    []database.Post
		NumPosts int
	}{
		Posts:    posts,
		NumPosts: numPosts,
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/blog.html")
	err = t.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
		sendResponseMsg("Failed to execute template", Error, w)
	}
}
