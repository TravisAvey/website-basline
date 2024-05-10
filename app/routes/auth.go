package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/travisavey/baseline/app/auth"
)

var (
	SESSION_NAME  string = "auth_session"
	AUTH_KEY      string = "authenticated"
	USER_ID       string = "user_id"
	AUTH_TOKEN    string = "auth_token"
	REFRESH_TOKEN string = "refresh_token"
	EXPIRES_IN    string = "expires_in"
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
		// TODO: log error
		w.Write([]byte(err.Error()))
	}
}

func authMiddleware(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetNamed(r, SESSION_NAME)
		if err != nil {
			// TODO: log error
			fmt.Println("authMiddleware session:", err.Error())
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

	session, err := auth.GetNamed(r, SESSION_NAME)
	if err != nil {
		// TODO: log error
		fmt.Println("loginAttempt session error:", err.Error())
	}

	session.Values[AUTH_KEY] = true
	session.Values[USER_ID] = user.User.ID
	session.Values[AUTH_TOKEN] = user.AccessToken
	session.Values[REFRESH_TOKEN] = user.RefreshToken
	session.Values[EXPIRES_IN] = user.ExpiresIn

	err = session.Save(r, w)
	if err != nil {
		// TODO: log error
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		return
	}

	w.Header().Add("Hx-Push-Url", "/dashboard")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func logOut(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetNamed(r, SESSION_NAME)
	if err != nil {
		// TODO: log error
		fmt.Println("logOut session", err.Error())
	}

	authToken := session.Values[AUTH_TOKEN]

	if authToken != nil {
		auth.SignOut(authToken.(string))
	}
	// incase we have no authToken
	// we can at least clear out the cookie
	// and the session
	session.Options.MaxAge = -1
	session.Save(r, w)

	c := &http.Cookie{
		Name:     SESSION_NAME,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	w.Header().Add("Hx-Push-Url", "/login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// TODO: doesn't show up after redirect
	msg := getResponseMsg("You have been logged out successfully", Success)
	sendLoginMessage(msg)
}
