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
	USER_ID      string = "user"
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
		w.WriteHeader(http.StatusUnauthorized)
		msg := getResponseMsg("incorrect username and/or password", Error)
		sendLoginMessage(msg)
		// TODO: log attempt
		fmt.Println(err.Error())
		return
	}
	// 1. store in session user logged in
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		// TODO: log error
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		return
	}
	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 12,
		HttpOnly: true,
		Path:     "/dashboard",
	}
	session.Values[AUTH_KEY] = true
	session.Values[USER_ID] = user.User.ID
	err = session.Store().Save(r, w, session)
	if err != nil {
		// TODO: log error
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		return
	}

	// 2. send user to dashboard?
	//
	w.Header().Add("Hx-Push-Url", "/dashboard")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
