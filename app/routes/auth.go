package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
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
		w.Write([]byte(err.Error()))
	}
}

func authMiddleware(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := auth.GetNamed(r, SESSION_NAME)

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

	session, _ := auth.GetNamed(r, SESSION_NAME)

	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 12,
		HttpOnly: true,
		Path:     "/dashboard",
	}

	fmt.Println("user.AccessToken...", user.AccessToken)
	session.Values[AUTH_KEY] = true
	session.Values[USER_ID] = user.User.ID
	session.Values[AUTH_TOKEN] = user.AccessToken
	session.Values[REFRESH_TOKEN] = user.RefreshToken
	session.Values[EXPIRES_IN] = user.ExpiresIn
	fmt.Println("Auth Token...", session.Values[AUTH_TOKEN])
	fmt.Println("session values: ", session.Values)

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
	fmt.Println(r.Cookies())
	session, _ := auth.GetNamed(r, SESSION_NAME)

	authToken := session.Values[AUTH_TOKEN]
	// authToken is always nil here =(
	fmt.Println("logOut: session.Value...", session.Values)
	fmt.Println("logOut: Auth Token...", authToken)

	if authToken != nil {
		auth.SignOut(authToken.(string))
	}

	session.Values[AUTH_KEY] = nil
	session.Values[USER_ID] = nil

	w.Header().Add("Hx-Push-Url", "/login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// TODO: doesn't show up after redirect
	msg := getResponseMsg("You have been logged out successfully", Success)
	sendLoginMessage(msg)
}
