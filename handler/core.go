package handler

import "github.com/anthdm/slick"

// RouteRegistrar is an interface that defines the method that a route registrar must implement
type RouteRegistrar interface {
	RegisterRoutes(app *slick.Slick)
}

// CompleteResourceController is an interface that defines the methods that a handler must implement
type CompleteResourceController interface {
	index(c *slick.Context) error
	show(c *slick.Context) error
	new(c *slick.Context) error
	create(c *slick.Context) error
	edit(c *slick.Context) error
	update(c *slick.Context) error
	delete(c *slick.Context) error
}
