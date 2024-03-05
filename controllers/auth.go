package controllers

import (
	"errors"
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
func (ctrl *AuthController) RegisterRoutes(app *fiber.App) {
	app.Get("/users", ctrl.index)
	app.Get("/users/:id", ctrl.show)
	app.Get("/profile", ctrl.profile)
	app.Get("/login", ctrl.login)
	app.Post("/login", ctrl.doLogin)
	app.Get("/signup", ctrl.signup)
	app.Post("/signup", ctrl.doSignup)
	app.Get("/profile/edit", ctrl.editProfile)
	app.Put("/profile", ctrl.updateProfile)
	app.Get("/forgot-password", ctrl.forgotPassword)
	app.Post("/forgot-password", ctrl.doForgotPassword)
	app.Get("/reset-password", ctrl.resetPassword)
	app.Put("/reset-password", ctrl.doResetPassword)
	app.Get("/logout", ctrl.logout)
}

// GET /users/ - index - List all users
func (ctrl *AuthController) index(ctx *fiber.Ctx) error {
	user := models.UserModel{}
	users, err := user.ReadAll(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("No users found"))
	}
	return RenderTempl(ctx, auth_views.Index(users))
}

// GET /users/:id - show - Show ctrl single user
func (ctrl *AuthController) show(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Invalid ID format"))
	}
	user := models.UserModel{ID: uint(id)}
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}
	return RenderTempl(ctx, auth_views.Show(user))
}

// GET /profile - profile - Show the profile of the logged in user
func (ctrl *AuthController) profile(ctx *fiber.Ctx) error {
	// get user from session
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}

	// get user from database
	user := models.UserModel{ID: userID}
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// render the profile
	return RenderTempl(ctx, auth_views.Profile(user))
}

// GET /login - login - Show the login form
func (ctrl *AuthController) login(ctx *fiber.Ctx) error {
	return RenderTempl(ctx, auth_views.Login())
}

// POST /login - doLogin - Process the login form
func (ctrl *AuthController) doLogin(ctx *fiber.Ctx) error {
	// get user from session
	_, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err == nil { // If no error, user is already logged in
		return ctx.Redirect("/profile")
	}

	// get user from database
	user := models.UserModel{Email: ctx.FormValue("email")}
	err = user.ReadByEmail(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ctx.FormValue("password")))
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Invalid password"))
	}

	// set user id in session
	err = SetUserInSession(ctx, ctrl.SessionStore, user.ID)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}

	// redirect to profile
	return ctx.Redirect("/profile")
}

// GET /signup - signup - Show the signup form
func (ctrl *AuthController) signup(ctx *fiber.Ctx) error {
	return RenderTempl(ctx, auth_views.Signup())
}

// POST /signup - doSignup - Process the signup form
func (ctrl *AuthController) doSignup(ctx *fiber.Ctx) error {
	// get user from session
	_, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err == nil { // If no error, user is already logged in
		return ctx.Redirect("/profile")
	}

	// create ctrl new user
	user := models.UserModel{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	// save the user
	err = user.Create(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to create user: "+err.Error()))
	}

	// set user id in session
	err = SetUserInSession(ctx, ctrl.SessionStore, user.ID)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}

	// redirect to profile
	return ctx.Redirect("/profile")
}

// GET /profile/edit - editProfile - Show the form to edit the profile of the logged in user
func (ctrl *AuthController) editProfile(ctx *fiber.Ctx) error {
	// get user from session
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}

	// get user from database
	user := models.UserModel{ID: userID}
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// render the edit profile form
	return RenderTempl(ctx, auth_views.EditProfile(user))
}

// PUT /profile - updateProfile - Update the profile of the logged in user
func (ctrl *AuthController) updateProfile(ctx *fiber.Ctx) error {
	// get user from session
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("You are not logged in"))
	}

	// get user from database
	user := models.UserModel{ID: userID}
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// update user
	if ctx.FormValue("email") != "" {
		user.Email = ctx.FormValue("email")
	}
	if ctx.FormValue("password") != "" {
		user.Password = ctx.FormValue("password")
	}

	// save the user
	err = user.Update(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to update user: "+err.Error()))
	}

	// redirect to profile
	return ctx.Redirect("/profile")
}

// GET /forgot-password - forgotPassword - Show the forgot password form
func (ctrl *AuthController) forgotPassword(ctx *fiber.Ctx) error {
	return RenderTempl(ctx, auth_views.ForgotPassword())
}

// POST /forgot-password - doForgotPassword - Process the forgot password form
func (ctrl *AuthController) doForgotPassword(ctx *fiber.Ctx) error {
	// get user from database
	user := models.UserModel{Email: ctx.FormValue("email")}
	err := user.ReadByEmail(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// create ctrl token
	token := models.TokenModel{
		Email: user.Email,
	}
	tokenValue, err := token.Create(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to create token: "+err.Error()))
	}

	// send email with token
	link := ctrl.Environment.BaseUrl + "/reset-password?email=" + user.Email + "&token=" + tokenValue
	ej := jobs.EmailJob{
		From:    "noreply@GoOnRails.com",
		To:      user.Email,
		Subject: "Reset your password",
		Body:    "Click this link to reset your password: " + link,
		Client:  ctrl.MailjetClient,
	}
	jobs.AddToQueue(ej)

	// TODO render and info message like "An email with instructions has been sent to your email address" instead
	// redirect to reset password
	return ctx.Redirect("/reset-password")
}

// GET /reset-password - resetPassword - Show the reset password form
func (ctrl *AuthController) resetPassword(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	tokenValue := ctx.Query("token")
	return RenderTempl(ctx, auth_views.ResetPassword(email, tokenValue))
}

// PUT /reset-password - doResetPassword - Process the reset password form
func (ctrl *AuthController) doResetPassword(ctx *fiber.Ctx) error {
	// get user from database
	user := models.UserModel{Email: ctx.FormValue("email")}
	err := user.ReadByEmail(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// get token from database
	token := models.TokenModel{
		Email: user.Email,
	}
	err = token.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("No valid tokens were found for this email, please request ctrl new one"))
	}

	// validate token value
	hashedToken, err := models.Hash(ctx.FormValue("token"))
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to hash input token: "+err.Error()))
	}
	if hashedToken != token.Value {
		return RenderTempl(ctx, auth_views.Error("Invalid token"))
	}

	// update user password
	user.Password = ctx.FormValue("password")
	err = user.Update(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to update user: "+err.Error()))
	}

	// delete token
	err = token.Delete(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to delete token: "+err.Error()))
	}

	// redirect to login
	return ctx.Redirect("/login")
}

// GET /logout - logout - Log the user out
func (ctrl *AuthController) logout(ctx *fiber.Ctx) error {
	// get user from session
	_, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Session error"))
	}

	// remove user id from session
	err = DeleteUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Session error"))
	}

	// redirect to login
	return ctx.Redirect("/login")
}

// Helper function to get the user from the session, if the user is not logged in, an error is returned.
func GetUserFromSession(ctx *fiber.Ctx, sessionStore *session.Store) (uint, error) {
	// get session
	sess, err := sessionStore.Get(ctx)
	if err != nil {
		return 0, errors.New("session error")
	}

	// get user id from session
	userID := sess.Get("user_id")
	if userID == nil {
		return 0, errors.New("you are not logged in")
	}
	return userID.(uint), nil
}

// Helper function to set the user in the session
func SetUserInSession(ctx *fiber.Ctx, sessionStore *session.Store, userID uint) error {
	// get session
	sess, err := sessionStore.Get(ctx)
	if err != nil {
		return errors.New("session error")
	}

	// set user id in session
	sess.Set("user_id", userID)
	err = sess.Save()
	if err != nil {
		return errors.New("session error")
	}
	return nil
}

// Helper function to delete the user from the session
func DeleteUserFromSession(ctx *fiber.Ctx, sessionStore *session.Store) error {
	// get session
	sess, err := sessionStore.Get(ctx)
	if err != nil {
		return errors.New("session error")
	}

	// remove user id from session
	sess.Delete("user_id")
	err = sess.Save()
	if err != nil {
		return errors.New("session error")
	}
	return nil
}
