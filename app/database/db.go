package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type conn_info struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
}

func getConnectionInfo() (string, error) {
	err := godotenv.Load("config/.env")
	if err != nil {
		return "", err
	}

	connInfo := conn_info{
		Host:     os.Getenv("POSTGRES_URL"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		connInfo.Host, connInfo.Port, connInfo.User, connInfo.Password, connInfo.DbName), nil
}

func Setup() error {
	dsn, err := getConnectionInfo()
	if err != nil {
		return err
	}

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
