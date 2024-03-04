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
	app.Get("/signup", a.signup)
	app.Post("/signup", a.doSignup)
	app.Get("/profile/edit", a.editProfile)
	app.Put("/profile", a.updateProfile)
	app.Get("/forgot-password", a.forgotPassword)
	app.Post("/forgot-password", a.doForgotPassword)
	app.Get("/logout", a.logout)
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
	user := models.User{Email: c.FormValue("email")}
	err = user.ReadByEmail()
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

// GET /signup - signup - Show the signup form
func (a *AuthController) signup(c *fiber.Ctx) error {
	return RenderTempl(c, auth.Signup())
}

// POST /signup - doSignup - Process the signup form
func (a *AuthController) doSignup(c *fiber.Ctx) error {
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

	// create a new user
	user := models.User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	// save the user
	err = user.Create()
	if err != nil {
		return RenderTempl(c, auth.Error("Failed to create user: "+err.Error()))
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

// GET /profile/edit - editProfile - Show the form to edit the profile of the logged in user
func (a *AuthController) editProfile(c *fiber.Ctx) error {
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

	// render the edit profile form
	return RenderTempl(c, auth.EditProfile(user))
}

// PUT /profile - updateProfile - Update the profile of the logged in user
func (a *AuthController) updateProfile(c *fiber.Ctx) error {
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

	// update user
	if c.FormValue("email") != "" {
		user.Email = c.FormValue("email")
	}
	if c.FormValue("password") != "" {
		user.Password = c.FormValue("password")
	}

	// save the user
	err = user.Update()
	if err != nil {
		return RenderTempl(c, auth.Error("Failed to update user: "+err.Error()))
	}

	// redirect to profile
	return c.Redirect("/profile")
}

// GET /forgot-password - forgotPassword - Show the forgot password form
func (a *AuthController) forgotPassword(c *fiber.Ctx) error {
	return RenderTempl(c, auth.ForgotPassword())
}

// POST /forgot-password - doForgotPassword - Process the forgot password form
func (a *AuthController) doForgotPassword(c *fiber.Ctx) error {
	// get user from database
	user := models.User{Email: c.FormValue("email")}
	err := user.ReadByEmail()
	if err != nil {
		return RenderTempl(c, auth.Error("User not found"))
	}

	// TODO send email with reset link

	// redirect to login
	return c.Redirect("/login")
}

// GET /logout - logout - Log the user out
func (a *AuthController) logout(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth.Error("Session error"))
	}

	// remove user id from session
	sess.Delete("user_id")
	err = sess.Save()
	if err != nil {
		return RenderTempl(c, auth.Error("Session error"))
	}

	// redirect to login
	return c.Redirect("/login")
}
