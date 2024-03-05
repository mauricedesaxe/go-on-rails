package controllers

import (
	hello_views "github.com/mauricedesaxe/go-on-rails/views/hello"

	"github.com/gofiber/fiber/v2"
)

type HelloController struct{}

func (h *HelloController) RegisterRoutes(app *fiber.App) {
	app.Get("/", h.index)
	app.Get("/hello", h.index)
}

// GET /posts - index - List all posts
func (h *HelloController) index(c *fiber.Ctx) error {
	return RenderTempl(c, hello_views.Index())
}
