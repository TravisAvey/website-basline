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
