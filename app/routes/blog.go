package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/travisavey/baseline/app/database"
)

func blog(w http.ResponseWriter, _ *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	t, _ := template.ParseFiles("web/templates/pages/blog.html")
	err = t.Execute(w, posts)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// create a post
func createPost(w http.ResponseWriter, r *http.Request) {
	post, err := parseFormData(r)
	if err != nil {
		w.Write([]byte("error parsing form"))
		return
	}

	err = database.NewPost(post)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// get a single post
func getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		// TODO: log Error
		w.Write([]byte(err.Error()))
		return
	}

	var post database.Post
	post, err = database.GetPost(id)
	if err != nil {
		// TODO: log error
		w.Write([]byte(err.Error()))
		return
	}

	t, _ := template.ParseFiles("web/templates/pages/blog/post.html")
	err = t.Execute(w, post)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func getPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	post, err := database.GetPostBySlug(slug)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/blog/post.html")

	t, _ := template.ParseFiles(files...)
	err = t.ExecuteTemplate(w, "base", post)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// get all posts
func getPosts(w http.ResponseWriter, _ *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		// TODO: log error
		w.Write([]byte("Error retreiving all posts"))
		return
	}

	// TODO: pass to template
	postStr := fmt.Sprint(posts)
	w.Write([]byte(postStr))
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	post, err := parseFormData(r)
	if err != nil {
		w.Write([]byte("error parsing form"))
		return
	}

	// errors updating:
	// pq: insert or update on table "post_categories" violates foreign key constraint "post_categories_post_id_fkey"
	// Thought I had this working before...
	err = database.UpdatePost(post)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// delete a post
func deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = database.DeletePost(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}
