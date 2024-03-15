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
	app.Get("/orders/new", ctrl.new)
	app.Post("/orders/validate", ctrl.validate)
	app.Post("/orders", ctrl.create)
	app.Get("/orders/:id", ctrl.show)
	app.Get("/orders/:id/edit", ctrl.edit)
	app.Put("/orders/:id", ctrl.update)
	app.Delete("/orders/:id", ctrl.delete)
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
	return RenderTempl(ctx, commerce_views.OrdersIndex(orders))
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
	return RenderTempl(ctx, commerce_views.OrdersShow(order))
}

// GET /orders/new - new - Show a form to create a new order
func (ctrl *OrdersController) new(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to create orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to create orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to create orders"))
	}

	// render a template with the new order form
	return RenderTempl(ctx, commerce_views.OrdersNew())
}

// POST /orders - create - Create a new order
func (ctrl *OrdersController) create(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to create orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to create orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to create orders"))
	}

	// create a new order
	var order model.Order
	err = ctx.BodyParser(&order)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error parsing order data"))
	}
	err = order.Create(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error creating order: "+err.Error()))
	}

	// redirect to the new order
	return ctx.Redirect("/orders/" + strconv.Itoa(int(order.ID)))
}

// GET /orders/:id/edit - edit - Show a form to edit an order
func (ctrl *OrdersController) edit(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to edit orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to edit orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to edit orders"))
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

	// render a template with the edit order form
	return RenderTempl(ctx, commerce_views.OrdersEdit(order))
}

// PUT /orders/:id - update - Update an order
func (ctrl *OrdersController) update(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to update orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to update orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to update orders"))
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

	// update the order
	err = ctx.BodyParser(&order)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error parsing order data"))
	}
	err = order.Update(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error updating order"))
	}

	// redirect to the updated order
	return ctx.Redirect("/orders/" + strconv.Itoa(int(order.ID)))
}

// GET /orders/validate - validate - Validate an order & return a partial view
func (ctrl *OrdersController) validate(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.ErrorAlert("You must be logged in to validate orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.ErrorAlert("You must be logged in to validate orders"))
	}

	// validate user id
	userIdStr := ctx.FormValue("user_id")
	if userIdStr != "" {
		var userToCheck model.User
		userIdInt, err := strconv.Atoi(userIdStr)
		if err != nil {
			return RenderTempl(ctx, common_views.ErrorAlert("Invalid user ID"))
		}
		userToCheck.ID = uint(userIdInt)
		err = userToCheck.Read(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, common_views.ErrorAlert("User not found"))
		}
	}

	// validate status
	status := ctx.FormValue("status")
	if status != "" {
		s := model.OrderStatus(status)
		switch s {
		case model.Waiting, model.Confirming, model.Confirmed, model.Sending, model.PartiallyPaid, model.Finished, model.Failed, model.Expired:
			// Status is valid
		default:
			return RenderTempl(ctx, common_views.ErrorAlert("Invalid order status"))
		}
	}

	// validate price currency
	priceCurrency := ctx.FormValue("price_currency")
	if priceCurrency != "" {
		if priceCurrency != "usd" {
			return RenderTempl(ctx, common_views.ErrorAlert("Invalid price currency"))
		}
	}

	// validate invoice URL
	invoiceUrl := ctx.FormValue("invoice_url")
	if invoiceUrl != "" {
		if invoiceUrl[:29] != "https://api.nowpayments.io/" {
			return RenderTempl(ctx, common_views.ErrorAlert("Invalid invoice URL"))
		}
	}

	// validate product id
	productIdStr := ctx.FormValue("product_id")
	if productIdStr != "" {
		var productToCheck model.Product
		productIdInt, err := strconv.Atoi(productIdStr)
		if err != nil {
			return RenderTempl(ctx, common_views.ErrorAlert("Invalid product ID"))
		}
		productToCheck.ID = uint(productIdInt)
		err = productToCheck.Read(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, common_views.ErrorAlert("Product not found"))
		}
	}

	return RenderTempl(ctx, common_views.InfoAlert("Order seems valid"))
}

// DELETE /orders/:id - destroy - Delete an order
func (ctrl *OrdersController) delete(ctx *fiber.Ctx) error {
	// check if user is logged in
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to delete orders"))
	}

	// get user from db & check if admin
	var user model.User
	user.ID = userId
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("You must be logged in to delete orders"))
	}
	if user.Role != "admin" {
		return RenderTempl(ctx, common_views.Error("You must be an admin to delete orders"))
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

	// delete the order
	err = order.Delete(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, common_views.Error("Error deleting order"))
	}

	// redirect to the orders index
	return ctx.Redirect("/orders")
}
