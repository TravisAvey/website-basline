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
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
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
	if post.Article.Updated {
		post.Article.UpdatedStr = parseDate(post.Article.DateUpdated.Time)
	}

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

// get a single post -- only use from Dashboard.
func getPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := getResponseMsg("Something went wrong getting the post. ID parsing", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	var post database.Post
	post, err = database.GetPostByID(id)
	if err != nil {
		msg := getResponseMsg("Something went wrong getting the post. DB retrieving", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}

	content := []byte(post.Article.Content)
	post.Article.HTML = template.HTML(content)

	post.Article.PostedStr = parseDate(post.Article.DatePosted.Time)
	if post.Article.Updated {
		post.Article.UpdatedStr = parseDate(post.Article.DateUpdated.Time)
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/post.html")
	err = t.Execute(w, post)
	if err != nil {
		// TODO: Log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
		w.Write([]byte(err.Error()))
	}
}

func editPostView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := getResponseMsg("Something went wrong getting the post. ID parsing", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	var post database.Post
	post, err = database.GetPostByID(id)
	if err != nil {
		msg := getResponseMsg("Something went wrong getting the post. DB retrieving", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}

	content := []byte(post.Article.Content)
	post.Article.HTML = template.HTML(content)

	post.Article.PostedStr = parseDate(post.Article.DatePosted.Time)
	if post.Article.Updated {
		post.Article.UpdatedStr = parseDate(post.Article.DateUpdated.Time)
	}

	cats, err := database.GetBlogCategories()
	if err != nil {
		// error getting Categories
		// TODO: log error
		fmt.Println(err.Error())
		return
	}
	data := struct {
		AllCats []postCats
		Post    database.Post
	}{
		Post:    post,
		AllCats: checkPostCategories(cats, post.Categories),
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/edit-post.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: Log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
		fmt.Println(err.Error())
	}
}

// a func just for all post categories and a flag for an active category
func checkPostCategories(allCats, cats []database.Category) []postCats {
	var categories []postCats

	for _, cat := range allCats {
		curCats := postCats{
			Category: cat.Category,
			Selected: false,
		}
		for _, c := range cats {
			if c.Category == cat.Category {
				curCats.Selected = true
			}
		}
		categories = append(categories, curCats)
	}

	return categories
}

// create a post
func createPost(w http.ResponseWriter, r *http.Request) {
	post, err := parsePostForm(r, true)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end.",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		fmt.Println("parsePostForm error", err.Error())
		return
	}

	err = database.NewPost(post)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end. We couldn't create a New Post",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		fmt.Println("parsePostForm error", err.Error())
	}
	msg := getResponseMsg("Post Created Successfully", Success)
	sendSSEMessage(msg)
}

func newPost(w http.ResponseWriter, r *http.Request) {
	cats, err := database.GetBlogCategories()
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end. We couldn't create a New Post",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		w.Write([]byte(err.Error()))
	}

	data := struct {
		Categories []database.Category
	}{
		Categories: cats,
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/new-post.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: Log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
		w.Write([]byte(err.Error()))
	}
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	post, err := parsePostForm(r, false)
	if err != nil {
		msg := errMsg{
			ErrorCode: 500,
			Message:   "Sorry, something went wrong on our end",
			Title:     "_Server Error",
			ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
		}
		sendErrorTemplate(msg, w)
		// TODO: Log error
		fmt.Println(err.Error())
		return
	}

	post.Article.Updated = true
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
		fmt.Println(err.Error())
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

func dashboardPosts(w http.ResponseWriter, _ *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		msg := getResponseMsg("There was an error retrieving the posts from the DB", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}

	for i := range posts {
		posts[i].Article.PostedStr = parseDate(posts[i].Article.DatePosted.Time)
		if posts[i].Article.Updated {
			posts[i].Article.UpdatedStr = parseDate(posts[i].Article.DateUpdated.Time)
		}
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
		// TODO: log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
	}
}

func dashboardPostCount(w http.ResponseWriter, _ *http.Request) {
	count, err := database.GetPostCount()
	if err != nil {
		msg := getResponseMsg("There was an error retrieving post count", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	w.Write([]byte(strconv.FormatUint(count, 10)))
}
