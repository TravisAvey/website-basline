package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var (
	supa  *supabase.Client
	store = sessions.NewCookieStore([]byte(GetSessionKey(32)))
)

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "session-name")
}

func GetNamed(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type User struct {
	AccessToken          string `json:"access_token"`
	TokenType            string `json:"token_type"`
	RefreshToken         string `json:"refresh_token"`
	ProviderToken        string `json:"provider_token"`
	ProviderRefreshToken string `json:"provider_refresh_token"`
	ExpiresIn            int    `json:"expires_in"`
}

func setup() {
	err := godotenv.Load("config/.env")
	if err != nil {
		// TODO: log error
		fmt.Println(err.Error())
	}

	supaUrl := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_KEY")

	supa = supabase.CreateClient(supaUrl, supaKey, true)
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

func Init() {
	setup()
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
