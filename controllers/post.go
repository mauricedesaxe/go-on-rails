package controllers

import (
	"net/http"
	"strconv"

	models "github.com/mauricedesaxe/go-on-rails/models"
	posts_views "github.com/mauricedesaxe/go-on-rails/views/posts"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type PostsController struct {
	Database *gorm.DB
} // Implements RouteRegistrar and CoreHandler

func (h *PostsController) RegisterRoutes(app *fiber.App) {
	app.Get("/posts", h.index)
	app.Get("/posts/:id", h.show)
	app.Get("/posts/new", h.new)
	app.Post("/posts", h.create)
	app.Get("/posts/:id/edit", h.edit)
	app.Put("/posts/:id", h.update)
	app.Delete("/posts/:id", h.delete)
}

// GET /posts - index - List all posts
func (h *PostsController) index(c *fiber.Ctx) error {
	post := models.PostModel{}
	postsRows, err := post.ReadAll(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("No posts found"))
	}
	return RenderTempl(c, posts_views.Index(postsRows))
}

// GET /posts/:id - show - Show a single post
func (h *PostsController) show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts_views.Error("Invalid ID format"))
	}

	post := models.PostModel{ID: uint(id)}
	err = post.Read(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("Post not found"))
	}

	return RenderTempl(c, posts_views.Show(post))
}

// GET /posts/new - new - Show the form to create a new post
func (h *PostsController) new(c *fiber.Ctx) error {
	return RenderTempl(c, posts_views.New())
}

// POST /posts - create - Create a new post
func (h *PostsController) create(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	// Ensure that the required fields are provided.
	if title == "" || content == "" || author == "" {
		return RenderTempl(c, posts_views.Error("Title, content, and author are required"))
	}

	// Create a new post.
	post := models.PostModel{
		Title:   title,
		Content: content,
		Author:  author,
	}
	err := post.Create(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("Failed to create post"))
	}

	// Redirect to the new post.
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (h *PostsController) edit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts_views.Error("Invalid ID format"))
	}

	post := models.PostModel{ID: uint(id)}
	err = post.Read(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("Post not found"))
	}

	return RenderTempl(c, posts_views.Edit(post))
}

// PUT /posts/:id - update - Update a post
func (h *PostsController) update(c *fiber.Ctx) error {
	// Parse the post ID from the request.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts_views.Error("Invalid ID format"))
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	post := models.PostModel{ID: uint(id)}

	// Update the post fields if they are provided.
	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	if author != "" {
		post.Author = author
	}

	// Save the updated post.
	err = post.Update(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("Failed to update post"))
	}

	// Redirect to the updated post.
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (h *PostsController) delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts_views.Error("Invalid ID format"))
	}

	post := models.PostModel{ID: uint(id)}
	err = post.Delete(h.Database)
	if err != nil {
		return RenderTempl(c, posts_views.Error("Failed to delete post"))
	}

	return c.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*PostsController)(nil)
var _ CompleteResourceController = (*PostsController)(nil)
