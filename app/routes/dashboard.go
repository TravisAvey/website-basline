package routes

// TODO: setup routes for the dashboard
// will need auth working for a backend

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/travisavey/baseline/app/database"
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	count, err := database.GetMessageCount(true)
	if err != nil {
		// TODO: log error
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
		sendResponseMsg("Failed to execute template", Error, w)
	}
}

// get a single post -- only use from Dashboard.
func getPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sendResponseMsg("Failed parse the ID for the blog post", Error, w)
		// TODO: log Error
		return
	}

	var post database.Post
	post, err = database.GetPostByID(id)
	if err != nil {
		sendResponseMsg("Failed to get the blog post", Error, w)
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
		w.Write([]byte(err.Error()))
	}
}

func dashboardGallery(w http.ResponseWriter, _ *http.Request) {
	imgs, err := database.GetAllImages()
	if err != nil {
		sendResponseMsg("Failed to get images", Error, w)
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
		w.Write([]byte(err.Error()))
		sendResponseMsg("Failed to execute template", Error, w)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := database.GetAllMessages()
	if err != nil {
		sendResponseMsg("Couldn't retrieve the messages", Error, w)
		// TODO: log Error
		return
	}

	// send partial
	println(msgs[0].Message)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sendResponseMsg("Couldn't parse the id of the message", Error, w)
		// TODO: log Error
		return
	}

	var message database.Message
	message, err = database.GetMessage(id)
	if err != nil {
		sendResponseMsg("Failed to get Message", Error, w)
		// TODO: log Error
		return

	}

	// send partial
	println(message.Message)
}

func messageRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sendResponseMsg("Couldn't parse the ID for the message", Error, w)
		// TODO: log Error
		return
	}

	err = database.MessageRead(id)
	if err != nil {
		sendResponseMsg("Failed to mark the message as read", Error, w)
		// TODO: log Error
	}
}

func messageDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sendResponseMsg("Couldn't parse the ID for the message", Error, w)
		// TODO: log Error
		return
	}

	err = database.DeleteMessage(id)
	if err != nil {
		sendResponseMsg("Couldn't Delete the message", Error, w)
		// TODO: log Error
	}
}
