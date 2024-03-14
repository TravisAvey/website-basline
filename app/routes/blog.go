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

	post, posts := posts[0], posts[1:]
	for i := range posts {
		posts[i].Article.PostedStr = parseDate(posts[i].Article.DatePosted.Time)
	}
	numPosts := len(posts)

	data := struct {
		ImageURL   string
		Title      string
		Summary    string
		Slug       string
		DatePosted string
		Posts      []database.Post
		NumPosts   int
	}{
		Posts:      posts,
		ImageURL:   post.Article.ImageURL,
		Title:      post.Article.Title,
		Summary:    post.Article.Summary,
		Slug:       post.Article.Slug,
		DatePosted: parseDate(post.Article.DatePosted.Time),
		NumPosts:   numPosts,
	}

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/blog.html")
	t, _ := template.ParseFiles(files...)
	err = t.ExecuteTemplate(w, "base", data)
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
func getPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		// TODO: log Error
		w.Write([]byte(err.Error()))
		return
	}

	var post database.Post
	post, err = database.GetPostByID(id)
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

	content := []byte(post.Article.Content)
	post.Article.HTML = template.HTML(content)

	post.Article.PostedStr = parseDate(post.Article.DatePosted.Time)

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
