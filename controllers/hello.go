package controllers

import (
	hello_views "github.com/mauricedesaxe/go-on-rails/views/hello"

	"github.com/gofiber/fiber/v2"
)

type HelloController struct{}

func (ctrl *HelloController) RegisterRoutes(app *fiber.App) {
	app.Get("/", ctrl.index)
	app.Get("/hello", ctrl.index)
}

// GET /posts - index - List all posts
func (ctrl *HelloController) index(ctx *fiber.Ctx) error {
	return RenderTempl(ctx, hello_views.Index())
}
