package main

import (
	"os"

	"github.com/travisavey/baseline/app/auth"
	"github.com/travisavey/baseline/app/database"
	"github.com/travisavey/baseline/app/logging"
	"github.com/travisavey/baseline/app/routes"
)

func main() {
	logging.Setup()

	var err error

	auth.Setup()
	user, authErr := auth.SignIn(os.Args[1], os.Args[2])
	if authErr != nil {
		println(authErr.Error())
	}

	println(user.User.ID)

	err = database.Init()
	if err != nil {
		panic(err)
	}

	routes.Init()
}
