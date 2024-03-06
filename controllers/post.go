package controllers

import (
	"net/http"
	"strconv"

	models "github.com/mauricedesaxe/go-on-rails/models"
	posts_views "github.com/mauricedesaxe/go-on-rails/views/posts"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type PostsController struct {
	Database     *gorm.DB
	SessionStore *session.Store
} // Implements RouteRegistrar and CoreHandler

func (ctrl *PostsController) RegisterRoutes(app *fiber.App) {
	app.Get("/posts", ctrl.index)
	app.Get("/posts/:id", ctrl.show)
	app.Get("/posts/new", ctrl.new)
	app.Post("/posts", ctrl.create)
	app.Get("/posts/:id/edit", ctrl.edit)
	app.Put("/posts/:id", ctrl.update)
	app.Delete("/posts/:id", ctrl.delete)
}

// GET /posts - index - List all posts
func (ctrl *PostsController) index(ctx *fiber.Ctx) error {
	post := models.PostModel{}
	postsRows, err := post.ReadAll(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("No posts found"))
	}
	return RenderTempl(ctx, posts_views.Index(postsRows))
}

// GET /posts/:id - show - Show a single post
func (ctrl *PostsController) show(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Invalid ID format"))
	}

	post := models.PostModel{ID: uint(id)}
	err = post.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Post not found"))
	}

	author := models.UserModel{ID: post.AuthorID}
	err = author.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Author not found"))
	}

	return RenderTempl(ctx, posts_views.Show(post, author))
}

// GET /posts/new - new - Show the form to create a new post
func (ctrl *PostsController) new(ctx *fiber.Ctx) error {
	return RenderTempl(ctx, posts_views.New())
}

// POST /posts - create - Create a new post
func (ctrl *PostsController) create(ctx *fiber.Ctx) error {
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	author := ctx.FormValue("author")

	// Ensure that the required fields are provided.
	if title == "" || content == "" || author == "" {
		return RenderTempl(ctx, posts_views.Error("Title, content, and author are required"))
	}

	// Get user from session
	userID, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("You must be logged in to create a post"))
	}

	// Create a new post.
	post := models.PostModel{
		Title:    title,
		Content:  content,
		AuthorID: userID,
	}
	err = post.Create(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Failed to create post"))
	}

	// Redirect to the new post.
	return ctx.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (ctrl *PostsController) edit(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Invalid ID format"))
	}

	// Get user from session
	_, err = GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("You must be logged in to edit a post"))
	}

	// Find the post by ID.
	post := models.PostModel{ID: uint(id)}
	err = post.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Post not found"))
	}

	return RenderTempl(ctx, posts_views.Edit(post))
}

// PUT /posts/:id - update - Update a post
func (ctrl *PostsController) update(ctx *fiber.Ctx) error {
	// Parse the post ID from the request.
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Invalid ID format"))
	}

	// Get user from session
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("You must be logged in to edit a post"))
	}

	// Find the post by ID.
	post := models.PostModel{ID: uint(id)}
	err = post.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Post not found"))
	}

	// Ensure that the user is the author of the post.
	if post.AuthorID != userId {
		return RenderTempl(ctx, posts_views.Error("You are not the author of this post"))
	}

	// Update the post fields locally if they are provided.
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}

	// Save the updated post to the database.
	err = post.Update(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Failed to update post"))
	}

	// Redirect to the updated post.
	return ctx.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (ctrl *PostsController) delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Invalid ID format"))
	}

	// Get user from session
	userId, err := GetUserFromSession(ctx, ctrl.SessionStore)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("You must be logged in to delete a post"))
	}

	// Find the post by ID.
	post := models.PostModel{ID: uint(id)}
	err = post.Read(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Post not found"))
	}

	// Ensure that the user is the author of the post.
	if post.AuthorID != userId {
		return RenderTempl(ctx, posts_views.Error("You are not the author of this post"))
	}

	// Delete the post from the database.
	err = post.Delete(ctrl.Database)
	if err != nil {
		return RenderTempl(ctx, posts_views.Error("Failed to delete post"))
	}

	return ctx.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*PostsController)(nil)
var _ CompleteResourceController = (*PostsController)(nil)
