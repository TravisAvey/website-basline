package main

import (
	"github.com/travisavey/baseline/app/auth"
	"github.com/travisavey/baseline/app/database"
	"github.com/travisavey/baseline/app/logging"
	"github.com/travisavey/baseline/app/routes"
)

func main() {
	logging.Setup()

	var err error

	auth.Setup()

	err = database.Setup()
	if err != nil {
		panic(err)
	}

	routes.Setup()
}
