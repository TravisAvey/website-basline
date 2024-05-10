package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var (
	supa  *supabase.Client
	store *sessions.CookieStore
)

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "session-name")
}

func GetNamed(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Setup() {
	err := godotenv.Load("config/.env")
	if err != nil {
		// TODO: log error
		fmt.Println(err.Error())
	}

	supaUrl := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_KEY")

	supa = supabase.CreateClient(supaUrl, supaKey, true)

	authKey := securecookie.GenerateRandomKey(64)
	encryptKey := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore(authKey, encryptKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func GetSessionKey(n int) string {
	b := make([]byte, n)
	for i := range b {
		v, err := rand.Int(rand.Reader, big.NewInt(32))
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		b[i] = letters[v.Int64()]
	}
	return string(b)
}

func SignIn(email, password string) (*supabase.AuthenticatedDetails, error) {
	ctx := context.Background()

	return supa.Auth.SignIn(ctx, supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
}

func SignOut(token string) error {
	ctx := context.Background()
	return supa.Auth.SignOut(ctx, token)
}
