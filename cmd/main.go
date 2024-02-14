package main

import (
	"goblog/handler"
	"goblog/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	model.DB = model.InitDB("store.db")
	app.Get("/", handler.HandleHelloIndex)

	// Create a slice of route registrars
	registrars := []handler.RouteRegistrar{
		&handler.Posts{},
	}

	// Register routes for each registrar
	for _, registrar := range registrars {
		registrar.RegisterRoutes(app)
	}

	log.Fatal(app.Listen(":3000"))
}
