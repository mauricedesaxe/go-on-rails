package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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

// ResponseHandler defines a function that handles the response in a specific format.
type ResponseHandler func(*fiber.Ctx) error

// HandleResponse decides whether to handle the response as JSON or as a Templ rendering,
// based on the "Accept" header of the request.
func HandleResponse(c *fiber.Ctx, handleJSON ResponseHandler, handleTempl ResponseHandler) error {
	if c.Get("Content-Type") == "application/json" {
		return handleJSON(c)
	}
	return handleTempl(c)
}

func RenderTempl(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
