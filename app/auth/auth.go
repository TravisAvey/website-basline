package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var supa *supabase.Client

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
