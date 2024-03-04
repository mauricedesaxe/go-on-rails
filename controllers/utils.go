package controllers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// RouteRegistrar is an interface that defines the method that a route registrar must implement
type RouteRegistrar interface {
	RegisterRoutes(app *fiber.App)
}

// CompleteResourceController is an interface that defines the methods that a controller must implement
type CompleteResourceController interface {
	index(c *fiber.Ctx) error
	show(c *fiber.Ctx) error
	new(c *fiber.Ctx) error
	create(c *fiber.Ctx) error
	edit(c *fiber.Ctx) error
	update(c *fiber.Ctx) error
	delete(c *fiber.Ctx) error
}

func RenderTempl(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
