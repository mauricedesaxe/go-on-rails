package main

import (
	"goblog/handler"
	"log"

	"github.com/anthdm/slick"
)

func main() {
	app := slick.New()
	app.Get("/", handler.HandleHelloIndex)
	log.Fatal(app.Start())
}
