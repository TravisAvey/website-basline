package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/travisavey/baseline/app/database"
)

type postCats struct {
	Category string
	Selected bool
}

func dashboard(w http.ResponseWriter, _ *http.Request) {
	count, err := database.GetMessageCount(true)
	if err != nil {
		// TODO: log error
		msg := getResponseMsg("There was an error retrieving from the DB", Error)
		sendSSEMessage(msg)
		w.Write([]byte(err.Error()))
	}
	data := struct {
		Title    string
		MsgCount uint64
	}{
		Title:    "Dashboard Page",
		MsgCount: count,
	}

	files := []string{"web/templates/base.html", "web/templates/pages/dashboard.html"}
	t, _ := template.ParseFiles(files...)
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		// TODO: log error
		w.Write([]byte(err.Error()))
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
	}
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
		if posts[i].Article.DateUpdated.Valid {
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
	if post.Article.DateUpdated.Valid {
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
	if post.Article.DateUpdated.Valid {
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
		Post    database.Post
		AllCats []postCats
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
	post, err := parsePostForm(r)
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

func dashboardGallery(w http.ResponseWriter, _ *http.Request) {
	imgs, err := database.GetAllImages()
	if err != nil {
		msg := getResponseMsg("Failed to get the images", Error)
		sendSSEMessage(msg)
		// TODO: log error
		return
	}

	numImgs := len(imgs)

	data := struct {
		Images    []database.Image
		NumImages int
	}{
		Images:    imgs,
		NumImages: numImgs,
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/gallery.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: log error
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := database.GetAllMessages()
	if err != nil {
		msg := getResponseMsg("Failed to get the messages", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	for i := range msgs {
		msgs[i].DateStr = parseDate(msgs[i].Sent.Time)
	}

	data := struct {
		Messages []database.Message
	}{
		Messages: msgs,
	}

	t, _ := template.ParseFiles("web/templates/pages/dashboard/messages.html")
	err = t.Execute(w, data)
	if err != nil {
		// TODO: log error
		sendResponseMsg("Failed to execute template", Error, w)
		msg := getResponseMsg("There was an error creating the page", Error)
		sendSSEMessage(msg)
	}
}

func getMessageCount(w http.ResponseWriter, _ *http.Request) {
	count, err := database.GetMessageCount(true)
	if err != nil {
		msg := getResponseMsg("Failed to get the message count", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	w.Write([]byte(strconv.FormatUint(count, 10)))
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := getResponseMsg("Failed to get the message. ID parsing", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	var message database.Message
	message, err = database.GetMessage(id)
	if err != nil {
		msg := getResponseMsg("Failed to get the message from the DB.", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return

	}

	// send partial
	println(message.Message)
}

func messageRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := getResponseMsg("Failed to parse the ID for the message", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	err = database.MessageRead(id)
	if err != nil {
		msg := getResponseMsg("Failed to mark the message as read", Error)
		sendSSEMessage(msg)
		// TODO: log Error
	}
}

func messageDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		msg := getResponseMsg("Failed to parse the ID for the message", Error)
		sendSSEMessage(msg)
		// TODO: log Error
		return
	}

	err = database.DeleteMessage(id)
	if err != nil {
		msg := getResponseMsg("Failed to delete the message from the DB", Error)
		sendSSEMessage(msg)
		// TODO: log Error
	}

	msg := getResponseMsg("The message has been deleted", Success)
	sendSSEMessage(msg)
}
