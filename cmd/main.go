package main

import (
	"goblog/handler"
	"goblog/model"
	"log"

	"github.com/anthdm/slick"
)

func main() {
	app := slick.New()
	model.InitDB()
	app.Get("/", handler.HandleHelloIndex)

	// Create a slice of route registrars
	registrars := []handler.RouteRegistrar{
		&handler.Posts{},
	}

	// Register routes for each registrar
	for _, registrar := range registrars {
		registrar.RegisterRoutes(app)
	}

	log.Fatal(app.Start())
}
