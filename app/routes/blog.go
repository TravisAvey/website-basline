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
		w.Write([]byte("error parsing id"))
		return
	}

	var post database.Post
	post, err = database.GetPost(id)
	if err != nil {
		// TODO: log error
		w.Write([]byte("Error getting post."))
		return
	}

	// TODO: pass to template
	postStr := fmt.Sprint(post)
	w.Write([]byte(postStr))
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
