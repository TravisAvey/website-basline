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
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
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
		// TODO: Log error
		w.Write([]byte(err.Error()))
	}
}

func getPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	post, err := database.GetPostBySlug(slug)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		return
	}

	content := []byte(post.Article.Content)
	post.Article.HTML = template.HTML(content)

	post.Article.PostedStr = parseDate(post.Article.DatePosted.Time)

	files := getBaseTemplates()
	files = append(files, "web/templates/pages/blog/post.html")

	t, _ := template.ParseFiles(files...)
	err = t.ExecuteTemplate(w, "base", post)
	if err != nil {
		// TODO: Log error
		w.Write([]byte(err.Error()))
	}
}

// get all posts
func getPosts(w http.ResponseWriter, _ *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		return
	}

	// TODO: pass to template
	postStr := fmt.Sprint(posts)
	w.Write([]byte(postStr))
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	post, err := parsePostForm(r)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		return
	}

	// errors updating:
	// pq: insert or update on table "post_categories" violates foreign key constraint "post_categories_post_id_fkey"
	// Thought I had this working before...
	err = database.UpdatePost(post)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end--Unable to update the post!",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
	}
}

// delete a post
func deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: log error
		return
	}
	err = database.DeletePost(id)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		return
	}
}

func getBlogCategories(w http.ResponseWriter, _ *http.Request) {
	cats, err := database.GetBlogCategories()
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end getting Blog Categories",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		return
	}

	fmt.Println("Blog Categories", cats)
}
