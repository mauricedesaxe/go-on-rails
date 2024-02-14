package handler

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseHandler defines a function that handles the response in a specific format.
type ResponseHandler func(*fiber.Ctx) error

// HandleResponse decides whether to handle the response as JSON or as a Templ rendering,
// based on the "Accept" header of the request.
func HandleResponse(c *fiber.Ctx, handleJSON ResponseHandler, handleTempl ResponseHandler) error {
	if c.Get("Accept") == "application/json" {
		return handleJSON(c)
	}
	return handleTempl(c)
}
