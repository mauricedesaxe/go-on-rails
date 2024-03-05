package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/mauricedesaxe/go-on-rails/env"
	"github.com/mauricedesaxe/go-on-rails/jobs"
	models "github.com/mauricedesaxe/go-on-rails/models"
	auth_views "github.com/mauricedesaxe/go-on-rails/views/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	Database      *gorm.DB
	SessionStore  *session.Store
	Environment   *env.Env
	MailjetClient *mailjet.Client
}

// RegisterRoutes registers the auth-related routes
func (a *AuthController) RegisterRoutes(app *fiber.App) {
	app.Get("/users", a.index)
	app.Get("/users/:id", a.show)
	app.Get("/profile", a.profile)
	app.Get("/login", a.login)
	app.Post("/login", a.doLogin)
	app.Get("/signup", a.signup)
	app.Post("/signup", a.doSignup)
	app.Get("/profile/edit", a.editProfile)
	app.Put("/profile", a.updateProfile)
	app.Get("/forgot-password", a.forgotPassword)
	app.Post("/forgot-password", a.doForgotPassword)
	app.Get("/reset-password", a.resetPassword)
	app.Put("/reset-password", a.doResetPassword)
	app.Get("/logout", a.logout)
}

// GET /users/ - index - List all users
func (a *AuthController) index(c *fiber.Ctx) error {
	user := models.UserModel{}
	users, err := user.ReadAll(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("No users found"))
	}
	return RenderTempl(c, auth_views.Index(users))
}

// GET /users/:id - show - Show a single user
func (a *AuthController) show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, auth_views.Error("Invalid ID format"))
	}
	user := models.UserModel{ID: uint(id)}
	err = user.Read(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}
	return RenderTempl(c, auth_views.Show(user))
}

// GET /profile - profile - Show the profile of the logged in user
func (a *AuthController) profile(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// get user id from session
	userID := sess.Get("user_id")
	if userID == nil {
		return RenderTempl(c, auth_views.Error("You are not logged in"))
	}

	// get user from database
	user := models.UserModel{ID: userID.(uint)}
	err = user.Read(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// render the profile
	return RenderTempl(c, auth_views.Profile(user))
}

// GET /login - login - Show the login form
func (a *AuthController) login(c *fiber.Ctx) error {
	return RenderTempl(c, auth_views.Login())
}

// POST /login - doLogin - Process the login form
func (a *AuthController) doLogin(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// check if user is already logged in
	userID := sess.Get("user_id")
	if userID != nil {
		return c.Redirect("/profile")
	}

	// get user from database
	user := models.UserModel{Email: c.FormValue("email")}
	err = user.ReadByEmail(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
	if err != nil {
		return RenderTempl(c, auth_views.Error("Invalid password"))
	}

	// set user id in session
	sess.Set("user_id", user.ID)
	err = sess.Save()
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// redirect to profile
	return c.Redirect("/profile")
}

// GET /signup - signup - Show the signup form
func (a *AuthController) signup(c *fiber.Ctx) error {
	return RenderTempl(c, auth_views.Signup())
}

// POST /signup - doSignup - Process the signup form
func (a *AuthController) doSignup(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// check if user is already logged in
	userID := sess.Get("user_id")
	if userID != nil {
		return c.Redirect("/profile")
	}

	// create a new user
	user := models.UserModel{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	// save the user
	err = user.Create(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to create user: "+err.Error()))
	}

	// set user id in session
	sess.Set("user_id", user.ID)
	err = sess.Save()
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// redirect to profile
	return c.Redirect("/profile")
}

// GET /profile/edit - editProfile - Show the form to edit the profile of the logged in user
func (a *AuthController) editProfile(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// get user id from session
	userID := sess.Get("user_id")
	if userID == nil {
		return RenderTempl(c, auth_views.Error("You are not logged in"))
	}

	// get user from database
	user := models.UserModel{ID: userID.(uint)}
	err = user.Read(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// render the edit profile form
	return RenderTempl(c, auth_views.EditProfile(user))
}

// PUT /profile - updateProfile - Update the profile of the logged in user
func (a *AuthController) updateProfile(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// get user id from session
	userID := sess.Get("user_id")
	if userID == nil {
		return RenderTempl(c, auth_views.Error("You are not logged in"))
	}

	// get user from database
	user := models.UserModel{ID: userID.(uint)}
	err = user.Read(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// update user
	if c.FormValue("email") != "" {
		user.Email = c.FormValue("email")
	}
	if c.FormValue("password") != "" {
		user.Password = c.FormValue("password")
	}

	// save the user
	err = user.Update(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to update user: "+err.Error()))
	}

	// redirect to profile
	return c.Redirect("/profile")
}

// GET /forgot-password - forgotPassword - Show the forgot password form
func (a *AuthController) forgotPassword(c *fiber.Ctx) error {
	return RenderTempl(c, auth_views.ForgotPassword())
}

// POST /forgot-password - doForgotPassword - Process the forgot password form
func (a *AuthController) doForgotPassword(c *fiber.Ctx) error {
	// get user from database
	user := models.UserModel{Email: c.FormValue("email")}
	err := user.ReadByEmail(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// create a token
	token := models.TokenModel{
		Email: user.Email,
	}
	tokenValue, err := token.Create(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to create token: "+err.Error()))
	}

	// send email with token
	link := a.Environment.BaseUrl + "/reset-password?email=" + user.Email + "&token=" + tokenValue
	ej := jobs.EmailJob{
		From:    "noreply@GoOnRails.com",
		To:      user.Email,
		Subject: "Reset your password",
		Body:    "Click this link to reset your password: " + link,
		Client:  a.MailjetClient,
	}
	jobs.AddToQueue(ej)

	// TODO render and info message like "An email with instructions has been sent to your email address" instead
	// redirect to reset password
	return c.Redirect("/reset-password")
}

// GET /reset-password - resetPassword - Show the reset password form
func (a *AuthController) resetPassword(c *fiber.Ctx) error {
	email := c.Query("email")
	tokenValue := c.Query("token")
	return RenderTempl(c, auth_views.ResetPassword(email, tokenValue))
}

// PUT /reset-password - doResetPassword - Process the reset password form
func (a *AuthController) doResetPassword(c *fiber.Ctx) error {
	// get user from database
	user := models.UserModel{Email: c.FormValue("email")}
	err := user.ReadByEmail(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("User not found"))
	}

	// get token from database
	token := models.TokenModel{
		Email: user.Email,
	}
	err = token.Read(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("No valid tokens were found for this email, please request a new one"))
	}

	// validate token value
	hashedToken, err := models.Hash(c.FormValue("token"))
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to hash input token: "+err.Error()))
	}
	if hashedToken != token.Value {
		return RenderTempl(c, auth_views.Error("Invalid token"))
	}

	// update user password
	user.Password = c.FormValue("password")
	err = user.Update(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to update user: "+err.Error()))
	}

	// delete token
	err = token.Delete(a.Database)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Failed to delete token: "+err.Error()))
	}

	// redirect to login
	return c.Redirect("/login")
}

// GET /logout - logout - Log the user out
func (a *AuthController) logout(c *fiber.Ctx) error {
	// get session
	sess, err := a.SessionStore.Get(c)
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// remove user id from session
	sess.Delete("user_id")
	err = sess.Save()
	if err != nil {
		return RenderTempl(c, auth_views.Error("Session error"))
	}

	// redirect to login
	return c.Redirect("/login")
}
