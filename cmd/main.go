package main

import (
	"fmt"

	"github.com/travisavey/baseline/app/database"
	"github.com/travisavey/baseline/app/model"
	"github.com/travisavey/baseline/app/routes"
)

func main() {
	fmt.Println("Hello!")
	model.Init()
	database.Init()
	routes.Init()
}
