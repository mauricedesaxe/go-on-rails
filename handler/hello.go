package handler

import (
	"goblog/view/hello"

	"github.com/gofiber/fiber/v2"
)

func HandleHelloIndex(c *fiber.Ctx) error {
	return RenderTempl(c, hello.Index())
}
