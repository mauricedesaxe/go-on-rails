package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	controllers "github.com/mauricedesaxe/go-on-rails/controllers"
	"github.com/mauricedesaxe/go-on-rails/jobs"
	models "github.com/mauricedesaxe/go-on-rails/models"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var sessionStore *session.Store

func init() {
	storage := sqlite3.New() // From github.com/gofiber/storage/sqlite3
	sessionStore = session.New(session.Config{
		Storage: storage,
	})
}

func main() {
	app := fiber.New()
	models.DB = models.InitDB("store.db")

	app.Use(fiberLogger.New())

	app.Static("/static", "./public")

	// Create a slice of route registrars
	registrars := []controllers.RouteRegistrar{
		&controllers.Hello{},
		&controllers.Posts{},
		&controllers.AuthController{SessionStore: sessionStore},
	}

	// Register routes for each registrar
	for _, registrar := range registrars {
		registrar.RegisterRoutes(app)
	}

	go jobs.StartSchedule()
	go jobs.StartQueue()

	// Graceful shutdown for job processing
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		jobs.StopQueue()
		log.Println("Graceful shutdown")
		os.Exit(0)
	}()

	log.Fatal(app.Listen(":3000"))
}
