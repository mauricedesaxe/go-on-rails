package controllers

import (
	"net/http"
	"strconv"

	models "github.com/mauricedesaxe/go-on-rails/models"
	"github.com/mauricedesaxe/go-on-rails/views/posts"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type Posts struct {
	Database *gorm.DB
} // Implements RouteRegistrar and CoreHandler

func (h *Posts) RegisterRoutes(app *fiber.App) {
	app.Get("/posts", h.index)
	app.Get("/posts/:id", h.show)
	app.Get("/posts/new", h.new)
	app.Post("/posts", h.create)
	app.Get("/posts/:id/edit", h.edit)
	app.Put("/posts/:id", h.update)
	app.Delete("/posts/:id", h.delete)
}

// GET /posts - index - List all posts
func (h *Posts) index(c *fiber.Ctx) error {
	post := models.Post{}
	postsRows, err := post.ReadAll(h.Database)
	if err != nil {
		return RenderTempl(c, posts.Error("No posts found"))
	}
	return RenderTempl(c, posts.Index(postsRows))
}

// GET /posts/:id - show - Show a single post
func (h *Posts) show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	post := models.Post{ID: uint(id)}
	err = post.Read(h.Database)
	if err != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	return RenderTempl(c, posts.Show(post))
}

// GET /posts/new - new - Show the form to create a new post
func (h *Posts) new(c *fiber.Ctx) error {
	return RenderTempl(c, posts.New())
}

// POST /posts - create - Create a new post
func (h *Posts) create(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	// Ensure that the required fields are provided.
	if title == "" || content == "" || author == "" {
		return RenderTempl(c, posts.Error("Title, content, and author are required"))
	}

	// Create a new post.
	post := models.Post{
		Title:   title,
		Content: content,
		Author:  author,
	}
	err := post.Create(h.Database)
	if err != nil {
		return RenderTempl(c, posts.Error("Failed to create post"))
	}

	// Redirect to the new post.
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (h *Posts) edit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	post := models.Post{ID: uint(id)}
	err = post.Read(h.Database)
	if err != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	return RenderTempl(c, posts.Edit(post))
}

// PUT /posts/:id - update - Update a post
func (h *Posts) update(c *fiber.Ctx) error {
	// Parse the post ID from the request.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	post := models.Post{ID: uint(id)}

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
		return RenderTempl(c, posts.Error("Failed to update post"))
	}

	// Redirect to the updated post.
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (h *Posts) delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	post := models.Post{ID: uint(id)}
	err = post.Delete(h.Database)
	if err != nil {
		return RenderTempl(c, posts.Error("Failed to delete post"))
	}

	return c.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*Posts)(nil)
var _ CompleteResourceController = (*Posts)(nil)
