package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var supa *supabase.Client

func getSecrets() {
	err := godotenv.Load("config/.env")
	if err != nil {
		// TODO: log error
		fmt.Println(err.Error())
	}

	supaUrl := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_KEY")

	supa = supabase.CreateClient(supaUrl, supaKey, true)
}

func Setup() {
	getSecrets()
}

func SignIn(email, password string) (*supabase.AuthenticatedDetails, error) {
	ctx := context.Background()

	return supa.Auth.SignIn(ctx, supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
}
