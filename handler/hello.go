package handler

import (
	"goblog/views/hello"

	"github.com/gofiber/fiber/v2"
)

func HandleHelloIndex(c *fiber.Ctx) error {
	return RenderTempl(c, hello.Index())
}
