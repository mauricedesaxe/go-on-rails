package handler

import (
	"errors"
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
	post, err := h.extractPostFromForm(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tx := model.DB.Create(&post)
	if tx.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.Status(http.StatusCreated).JSON(post)
}

// createTempl handles the Templ response for creating a new post.
func (h *Posts) createTempl(c *fiber.Ctx) error {
	post, err := h.extractPostFromForm(c)
	if err != nil {
		return RenderTempl(c, posts.Error(err.Error()))
	}

	tx := model.DB.Create(&post)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Failed to create post"))
	}

	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// extractPostFromForm extracts post data from the form values.
func (h *Posts) extractPostFromForm(c *fiber.Ctx) (model.Post, error) {
	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	if title == "" || content == "" || author == "" {
		return model.Post{}, errors.New("missing title, content, or author")
	}

	return model.Post{
		Title:   title,
		Content: content,
		Author:  author,
	}, nil
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	if err := h.extractAndUpdatePost(c, &res); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusAccepted).JSON(res)
}

// updateTempl handles the Templ rendering for updating a post.
func (h *Posts) updateTempl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return RenderTempl(c, posts.Error("Post not found"))
	}

	if err := h.extractAndUpdatePost(c, &res); err != nil {
		return RenderTempl(c, posts.Error(err.Error()))
	}

	return c.Redirect("/posts/"+strconv.Itoa(int(res.ID)), http.StatusFound)
}

// extractAndUpdatePost extracts post data from the form and updates the post.
func (h *Posts) extractAndUpdatePost(c *fiber.Ctx, post *model.Post) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

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

	tx := model.DB.Save(post)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
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
