package handler

import (
	"github.com/gofiber/fiber/v2"
)

// RouteRegistrar is an interface that defines the method that a route registrar must implement
type RouteRegistrar interface {
	RegisterRoutes(app *fiber.App)
}

// CompleteResourceController is an interface that defines the methods that a handler must implement
type CompleteResourceController interface {
	index(c *fiber.Ctx) error
	show(c *fiber.Ctx) error
	new(c *fiber.Ctx) error
	create(c *fiber.Ctx) error
	edit(c *fiber.Ctx) error
	update(c *fiber.Ctx) error
	delete(c *fiber.Ctx) error
}
