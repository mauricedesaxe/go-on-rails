package handler

import (
	"goblog/model"
	"goblog/view/posts"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Posts struct{} // Implements RouteRegistrar and CoreHandler

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
	return HandleResponse(c, h.indexJSON, h.indexTempl)
}

// indexJSON handles the JSON response for listing all posts.
func (h *Posts) indexJSON(c *fiber.Ctx) error {
	var rows []model.Post
	tx := model.DB.Find(&rows)
	if tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve posts"})
	}
	return c.JSON(rows)
}

// indexTempl handles the Templ rendering for listing all posts.
func (h *Posts) indexTempl(c *fiber.Ctx) error {
	var rows []model.Post
	tx := model.DB.Find(&rows)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("No posts found"))
	}
	return RenderTempl(c, posts.Index(rows))
}

// GET /posts/:id - show - Show a single post
func (h *Posts) show(c *fiber.Ctx) error {
	return HandleResponse(c, h.showJSON, h.showTempl)
}

// showJSON handles the JSON response for showing a single post.
func (h *Posts) showJSON(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	return c.JSON(res)
}

// showTempl handles the Templ rendering for showing a single post.
func (h *Posts) showTempl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	return RenderTempl(c, posts.Show(res))
}

// GET /posts/new - new - Show the form to create a new post
func (h *Posts) new(c *fiber.Ctx) error {
	return HandleResponse(c, h.newJSON, h.newTempl)
}

// newJSON handles the case where JSON response is requested for a route that does not support it.
func (h *Posts) newJSON(c *fiber.Ctx) error {
	return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{"error": "JSON not supported for creating new posts"})
}

// newTempl handles rendering the form to create a new post using Templ.
func (h *Posts) newTempl(c *fiber.Ctx) error {
	return RenderTempl(c, posts.New())
}

// POST /posts - create - Create a new post
func (h *Posts) create(c *fiber.Ctx) error {
	return HandleResponse(c, h.createJSON, h.createTempl)
}

// createJSON handles the JSON response for creating a new post.
func (h *Posts) createJSON(c *fiber.Ctx) error {
	var post model.Post

	// Parse the request body into the post model.
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Ensure that the required fields are provided.
	if post.Title == "" || post.Content == "" || post.Author == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Title, content, and author are required"})
	}

	// Save the new post.
	tx := model.DB.Create(&post)
	if tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	// Return the new post.
	return c.Status(http.StatusCreated).JSON(post)
}

// createTempl handles the Templ response for creating a new post.
func (h *Posts) createTempl(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	// Ensure that the required fields are provided.
	if title == "" || content == "" || author == "" {
		return RenderTempl(c, posts.Error("Title, content, and author are required"))
	}

	// Create a new post.
	post := model.Post{
		Title:   title,
		Content: content,
		Author:  author,
	}

	// Save the new post.
	tx := model.DB.Create(&post)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Failed to create post"))
	}

	// Redirect to the new post.
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (h *Posts) edit(c *fiber.Ctx) error {
	return HandleResponse(c, h.editJSON, h.editTempl)
}

// editJSON handles the JSON response for editing a post.
func (h *Posts) editJSON(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	return c.JSON(res)
}

// editTempl handles the Templ rendering for editing a post.
func (h *Posts) editTempl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	return RenderTempl(c, posts.Edit(res))
}

// PUT /posts/:id - update - Update a post
func (h *Posts) update(c *fiber.Ctx) error {
	return HandleResponse(c, h.updateJSON, h.updateTempl)
}

// updateJSON handles the JSON response for updating a post.
func (h *Posts) updateJSON(c *fiber.Ctx) error {
	// Parse the post ID from the request.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Find the post to update.
	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	// Parse the request body into the post model.
	var body model.Post
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// Update the post fields if they are provided.
	if body.Title != "" {
		res.Title = body.Title
	}
	if body.Content != "" {
		res.Content = body.Content
	}
	if body.Author != "" {
		res.Author = body.Author
	}

	// Save the updated post.
	tx = model.DB.Save(&res)
	if tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	// Return the updated post.
	return c.Status(http.StatusOK).JSON(res)
}

// updateTempl handles the Templ rendering for updating a post.
func (h *Posts) updateTempl(c *fiber.Ctx) error {
	// Parse the post ID from the request.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	// Find the post to update.
	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	// Update the post fields if they are provided.
	if title != "" {
		res.Title = title
	}
	if content != "" {
		res.Content = content
	}
	if author != "" {
		res.Author = author
	}

	// Save the updated post.
	tx = model.DB.Save(res)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Failed to update post"))
	}

	// Redirect to the updated post.
	return c.Redirect("/posts/"+strconv.Itoa(int(res.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (h *Posts) delete(c *fiber.Ctx) error {
	return HandleResponse(c, h.deleteJSON, h.deleteTempl)
}

// deleteJSON handles the JSON response for deleting a post.
func (h *Posts) deleteJSON(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	tx = model.DB.Delete(&res)
	if tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete post"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Post deleted"})
}

// deleteTempl handles the Templ rendering for deleting a post.
func (h *Posts) deleteTempl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	tx = model.DB.Delete(&res)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Failed to delete post"))
	}

	return c.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*Posts)(nil)
var _ CompleteResourceController = (*Posts)(nil)
