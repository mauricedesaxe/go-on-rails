package controllers

import (
	"goblog/views/hello"

	"github.com/gofiber/fiber/v2"
)

type Hello struct{}

func (h *Hello) RegisterRoutes(app *fiber.App) {
	app.Get("/", h.index)
	app.Get("/hello", h.index)
}

// GET /posts - index - List all posts
func (h *Hello) index(c *fiber.Ctx) error {
	return HandleResponse(c, h.indexJSON, h.indexTempl)
}

// indexJSON handles the JSON response for listing all posts.
func (h *Hello) indexJSON(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}

// indexTempl handles the Templ rendering for listing all posts.
func (h *Hello) indexTempl(c *fiber.Ctx) error {
	return RenderTempl(c, hello.Index())
}
