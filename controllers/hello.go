package controllers

import (
	"github.com/mauricedesaxe/go-on-rails/views/hello"

	"github.com/gofiber/fiber/v2"
)

type Hello struct{}

func (h *Hello) RegisterRoutes(app *fiber.App) {
	app.Get("/", h.index)
	app.Get("/hello", h.index)
}

// GET /posts - index - List all posts
func (h *Hello) index(c *fiber.Ctx) error {
	return RenderTempl(c, hello.Index())
}
