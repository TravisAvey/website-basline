package main

import (
	"github.com/travisavey/baseline/app/database"
	"github.com/travisavey/baseline/app/logging"
	"github.com/travisavey/baseline/app/model"
	"github.com/travisavey/baseline/app/routes"
)

func main() {
	var err error

	model.Init()

	err = database.Init()
	if err != nil {
		panic(err)
	}

	logging.Setup()

	routes.Init()
}
