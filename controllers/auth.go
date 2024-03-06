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
	// get user from database
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Invalid ID format"))
	}
	user := models.UserModel{ID: uint(id)}
	err = user.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("User not found"))
	}

	// get user from session to check if it's the logged in user
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}
	var isLoggedUser bool = false
	if userID != user.ID {
		isLoggedUser = true
	}

	// render the user
	return RenderTempl(ctx, auth_views.Show(user, isLoggedUser))
}

// GET /profile - profile - Redirect to the user's profile
func (ctrl *AuthController) profile(ctx *fiber.Ctx) error {
	// get user from session
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error(err.Error()))
	}

	return ctx.Redirect("/users/" + strconv.Itoa(int(userID)))
}

// GET /login - login - Show the login form
func (ctrl *AuthController) login(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	token := ctx.Query("token")

	// if email and token are present, authenticate the user
	if email != "" && token != "" {
		// get user from database
		user := models.UserModel{Email: email}
		err := user.ReadByEmail(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error("User not found"))
		}

		// get token from database
		tokenModel := models.TokenModel{Email: user.Email}
		err = tokenModel.Read(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error("No valid tokens were found for this email, please request ctrl new one"))
		}

		// validate token value
		hashedToken, err := models.Hash(token)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error("Failed to hash input token: "+err.Error()))
		}
		if hashedToken != tokenModel.Value {
			return RenderTempl(ctx, auth_views.Error("Invalid token"))
		}

		// delete token
		err = tokenModel.Delete(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error("Failed to delete token: "+err.Error()))
		}

		// set user id in session
		err = SetUserInSession(ctx, ctrl.SessionStore, user.ID)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error(err.Error()))
		}

		// redirect to profile
		return ctx.Redirect("/profile")
	}

	// if email and token are not present, show the login form
	return RenderTempl(ctx, auth_views.Login())
}

// POST /login - doLogin - Send the magic link
func (ctrl *AuthController) doLogin(ctx *fiber.Ctx) error {
	var isNewUser bool = false

	// get user from database
	user := models.UserModel{Email: ctx.FormValue("email")}
	err := user.ReadByEmail(ctrl.Database)
	if err != nil {
		// if user doesn't exist, create a new user
		err = user.Create(ctrl.Database)
		if err != nil {
			return RenderTempl(ctx, auth_views.Error("Failed to create user: "+err.Error()))
		}
		isNewUser = true
	}

	// create a new token
	token := models.TokenModel{
		Email: user.Email,
	}
	unhashedTokenValue, err := token.Create(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, auth_views.Error("Failed to create token: "+err.Error()))
	}

	// send magic link email
	link := ctrl.Environment.BaseUrl + "/magic-link?email=" + user.Email + "&token=" + unhashedTokenValue
	ej := jobs.EmailJob{
		From:    ctrl.Environment.FromEmail,
		To:      user.Email,
		Subject: "Log in to GoOnRails",
		Body:    "Click this link to log in: " + link,
		Client:  ctrl.MailjetClient,
	}
	jobs.AddToQueue(ej)

	infoMsg := "An email with login instructions has been sent to your email address"
	if isNewUser {
		infoMsg = "We created a new account for you. An email with login instructions has been sent to your email address."
	}
	return RenderTempl(ctx, auth_views.Info(infoMsg))
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
