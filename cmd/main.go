package main

import (
	"log"

	"github.com/dbo-test/internal/app"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	// Add Router
	app.HTTPServers.Start(":8080")
}
