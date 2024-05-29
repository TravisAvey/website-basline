package services

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

// client object to interface with the aws-sdk-go service
var s3Client *s3.S3

// the struct for the env vars
type S3Storage struct {
	Key      string
	Secret   string
	Endpoint string
	Bucket   string
	Url      string
}

// the storage object to house all the env vars
var s3Storage S3Storage

// InitS3Storage Initializes the S3 Storage objects and environmental variables
func InitS3Storage() {
	// Grab the env vars
	err := godotenv.Load("config/.env")
	if err != nil {
		fmt.Printf("Error loading the environment file: %v\n", err)
		// TODO: log error
		// should we panic.. or set flag that s3 storage not up.
	}
	// store the env vars
	s3Storage.Key = os.Getenv("S3_ACCESS_TOKEN")
	s3Storage.Secret = os.Getenv("S3_SECRET")
	s3Storage.Endpoint = os.Getenv("S3_ENDPOINT")
	s3Storage.Bucket = os.Getenv("S3_BUCKET")
	s3Storage.Url = os.Getenv("S3_URL")

	// init the aws-sdk
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(s3Storage.Key, s3Storage.Secret, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	// start the new sessions for the aws-sdk
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		// fatal or do we set a flag ~ s3 failed to init, so we pass on any func call..
		fmt.Println(err.Error())
		return
	}
	s3Client = s3.New(newSession)
}

// SendImage sends the image from the fileLocation string to s3 storage
// fileLocation the location on disk where the file is stored locally
// image the name of the file that will be in s3 storage -- a key, basically
//
//	where https://s3.storage.com/image.jpg, image.jpg will be the key
func SendImage(fileLocation string, image string) error {
	file, fileErr := os.Open(fileLocation)
	if fileErr != nil {
		fmt.Printf("error opening the file")
		return fileErr
	}
	object := s3.PutObjectInput{
		// bucket is space name
		Bucket: aws.String(s3Storage.Bucket),
		// key is folder/filename.ext
		Key:  aws.String(image),
		Body: file,
		ACL:  aws.String("public-read"),
	}
	_, err := s3Client.PutObject(&object)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// GetImage Returns the image file from the S3 Storage
// imageKey - the file from S3 Storage: image.jpg
// would download the image.jpg from S3 Storage and then return the file location
func GetImage(imageKey string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3Storage.Bucket),
		Key:    aws.String(imageKey),
	}

	result, err := s3Client.GetObject(input)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	file := fmt.Sprintf("/temp/%v", imageKey)
	out, err := os.Create(file)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err.Error())
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, result.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return file, nil
}

// DeleteImage removes the image from the S3 Storage
// image - the filename of the image: image.jpg, where the url of the image would be https://s3.storage.com/image.jpg
func DeleteImage(image string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Storage.Bucket),
		Key:    aws.String(image),
	}

	_, err := s3Client.DeleteObject((input))
	if err != nil {
		fmt.Printf("Error with deleting image from s3 storage: %v\n", err.Error())
		return err
	}
	return nil
}

// GetUrl returns the URL endpoint to be used with the image filename for the DB and image location
// this URL is the S3_URL from the environment file, will be updated per website's bucket in S3 storage
func GetS3Url() string {
	return s3Storage.Url
}

// GetS3Key returns the key/ID for the S3 Storage
func GetS3Key() string {
	return s3Storage.Key
}

// GetS3Endpoint returns the endpoint for the S3 Storage
// Not to be confused with the URL, this is not to be used for image location endpoints
// Used for the endpoint for the S3 Storage Endpoint Configuration
func GetS3Endpoint() string {
	return s3Storage.Endpoint
}

// GetS3Bucket returns the bucket for the S3 Storage.
// Could be useful for storing images for different parts of the site
// Will be different for each website
func GetS3Bucket() string {
	return s3Storage.Bucket
}
