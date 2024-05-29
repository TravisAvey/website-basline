package services

import (
	"fmt"
	"os"

	Tinify "github.com/gwpp/tinify-go/tinify"
	"github.com/joho/godotenv"
)

var apiKey string

func InitTinify() {
	err := godotenv.Load("config/.env")
	if err != nil {
		// TODO: log here? or panic?
		fmt.Printf("Error loading the environment file: %v\n", err)
		return
	}

	apiKey = os.Getenv("TINYPNG_KEY")
	fmt.Println("apiKey", apiKey)
}

// CompressImage compresses an image (file - location of image)
func CompressImage(file string, destination string) error {
	Tinify.SetKey(apiKey)
	source, err := Tinify.FromFile(file)
	if err != nil {
		return err
	}

	err = source.ToFile(destination)
	if err != nil {
		return err
	}
	return nil
}
