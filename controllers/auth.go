package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	models "github.com/mauricedesaxe/go-on-rails/models"
	"github.com/mauricedesaxe/go-on-rails/views/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	SessionStore *session.Store
}

// RegisterRoutes registers the auth-related routes
func (a *AuthController) RegisterRoutes(app *fiber.App) {
	app.Get("/profile", a.profile)
	app.Get("/login", a.login)
	app.Post("/login", a.doLogin)
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

// POST /login - doLogin - Process the login form
func (a *AuthController) doLogin(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth.Error("Session error"))
	}

	// check if user is already logged in
	userID := sess.Get("user_id")
	if userID != nil {
		return c.Redirect("/profile")
	}

	// get user from database
	user := models.User{Username: c.FormValue("username")}
	err = user.ReadByUsername()
	if err != nil {
		return RenderTempl(c, auth.Error("User not found"))
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
	if err != nil {
		return RenderTempl(c, auth.Error("Invalid password"))
	}

	// set user id in session
	sess.Set("user_id", user.ID)
	err = sess.Save()
	if err != nil {
		return RenderTempl(c, auth.Error("Session error"))
	}

	// redirect to profile
	return c.Redirect("/profile")
}
