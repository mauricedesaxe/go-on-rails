package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mauricedesaxe/go-on-rails/env"
	model "github.com/mauricedesaxe/go-on-rails/models"
	commerce_views "github.com/mauricedesaxe/go-on-rails/views/commerce"
	common_views "github.com/mauricedesaxe/go-on-rails/views/common"
	"gorm.io/gorm"
)

type OrdersController struct {
	Database     *gorm.DB
	SessionStore *session.Store
	Environment  *env.Env
}

func (ctrl *OrdersController) RegisterRoutes(app *fiber.App) {
	app.Get("/orders", ctrl.index)
	app.Get("/orders/:id", ctrl.show)
}

// GET /orders - index - List all orders
func (ctrl *OrdersController) index(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to view orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to view orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to view orders"))
	}

	// get all orders from db
	var order model.Order
	var orders []model.Order
	orders, err = order.ReadAll(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error reading orders from database"))
	}

	// render a template with the orders
	return RenderTempl(ctx, commerce_views.Index(orders))
}

// GET /orders/:id - show - Show a single order
func (ctrl *OrdersController) show(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to view orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to view orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to view orders"))
	}

	// get order from db
	var order model.Order
	strID := ctx.Params("id")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Invalid order ID"))
	}
	order.ID = uint(intID)
	err = order.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error reading order from database"))
	}

	// render a template with the order
	return RenderTempl(ctx, commerce_views.Show(order))
}
