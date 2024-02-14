package handler

import (
	"goblog/view/hello"

	"github.com/anthdm/slick"
)

func HandleHelloIndex(c *slick.Context) error {
	return c.Render(hello.Index())
}
