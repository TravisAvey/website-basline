package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/travisavey/baseline/app/auth"
)

var (
	SESSION_NAME string = "auth_session"
	AUTH_KEY     string = "authenticated"
	USER_ID      string = "user_id"
	store               = sessions.NewCookieStore([]byte(auth.GetSessionKey(32)))
)

func loginPage(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title string
	}{
		Title: "_Login",
	}
	files := []string{"web/templates/base.html", "web/templates/pages/login.html"}
	t, _ := template.ParseFiles(files...)
	err := t.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func authMiddleware(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, SESSION_NAME)
		if err != nil {
			// TODO: log error
			msg := errMsg{
				ErrorCode: 500,
				Message:   "Sorry, something went wrong on our end",
				Title:     "_Server Error",
				ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
			}
			sendErrorTemplate(msg, w)
		}
		if session.Values[AUTH_KEY] == nil || session.Values[USER_ID] == nil {
			// user not auth...
			msg := errMsg{
				ErrorCode: 404,
				Message:   "You are not authorized to see this page. Please login first",
				Title:     "_Not Authorized",
				ImageURL:  "https://picsum.photos/1920/1080/?blur=2",
			}
			sendErrorTemplate(msg, w)
		}
		HandlerFunc.ServeHTTP(w, r)
	}
}

func loginAttempt(w http.ResponseWriter, r *http.Request) {
	user, err := auth.SignIn(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		// send msg that incorrect user/password
		sendResponseMsg("Incorrect username and/or password", Error, w)
		// TODO: log attempt
		fmt.Println(err.Error())
		return
	}
	// 1. store in session user logged in
	// 2. send user to dashboard?
	fmt.Println(user.User)
}
