package main

import (
	"goblog/handler"
	models "goblog/models"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	models.DB = models.InitDB("store.db")

	app.Use(fiberLogger.New())

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
