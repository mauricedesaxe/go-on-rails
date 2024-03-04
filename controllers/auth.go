package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	models "github.com/mauricedesaxe/go-on-rails/models"
	"github.com/mauricedesaxe/go-on-rails/views/auth"
)

type AuthController struct {
	SessionStore *session.Store
}

// RegisterRoutes registers the auth-related routes
func (a *AuthController) RegisterRoutes(app *fiber.App) {
	app.Get("/profile", a.profile)
	app.Get("/login", a.login)
}

// GET /profile - profile - Show the profile of the logged in user
func (a *AuthController) profile(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth.Error("Session error"))
	}

	// get user id from session
	userID := sess.Get("user_id")
	if userID == nil {
		return RenderTempl(c, auth.NotLoggedIn())
	}

	// get user from database
	user := models.User{ID: userID.(uint)}
	err = user.Read()
	if err != nil {
		return RenderTempl(c, auth.Error("User not found"))
	}

	// render the profile
	return RenderTempl(c, auth.Profile(user))
}

// GET /login - login - Show the login form
func (a *AuthController) login(c *fiber.Ctx) error {
	return RenderTempl(c, auth.Login())
}
