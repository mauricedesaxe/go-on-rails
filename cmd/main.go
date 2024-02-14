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
	log.Fatal(app.Start())
}
