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
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	var rows []model.Post
	tx := model.DB.Find(&rows)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve posts"})
		}
		return RenderTempl(c, posts.Error("No posts found"))
	}

	if isJSON {
		return c.JSON(rows)
	}
	return RenderTempl(c, posts.Index(rows))
}

// GET /posts/:id - show - Show a single post
func (h *Posts) show(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if isJSON {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
		}
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		return RenderTempl(c, posts.Error("Post not found"))
	}

	if isJSON {
		return c.JSON(res)
	}
	return RenderTempl(c, posts.Show(res))
}

// GET /posts/new - new - Show the form to create a new post
func (h *Posts) new(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	if isJSON {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{"error": "JSON not supported"})
	}
	return RenderTempl(c, posts.New())
}

// POST /posts - create - Create a new post
func (h *Posts) create(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")

	post := model.Post{
		Title:   title,
		Content: content,
		Author:  author,
	}

	tx := model.DB.Create(&post)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
		}
		return RenderTempl(c, posts.Error("Failed to create post"))
	}

	if isJSON {
		return c.Status(http.StatusCreated).JSON(post)
	}
	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (h *Posts) edit(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if isJSON {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
		}
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		return RenderTempl(c, posts.Error("Post not found"))
	}

	if isJSON {
		return c.JSON(res)
	}
	return RenderTempl(c, posts.Edit(res))
}

// PUT /posts/:id - update - Update a post
func (h *Posts) update(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if isJSON {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
		}
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		return RenderTempl(c, posts.Error("Post not found"))
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")
	res.Title = title
	res.Content = content
	res.Author = author
	tx = model.DB.Save(&res)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
		}
		return RenderTempl(c, posts.Error("Failed to update post"))
	}

	if isJSON {
		return c.JSON(res)
	}
	return c.Redirect("/posts/"+strconv.Itoa(int(res.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (h *Posts) delete(c *fiber.Ctx) error {
	var isJSON bool
	if c.Get("Accept") == "application/json" {
		isJSON = true
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if isJSON {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
		}
		return RenderTempl(c, posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		return RenderTempl(c, posts.Error("Post not found"))
	}

	tx = model.DB.Delete(&res)
	if tx.Error != nil {
		if isJSON {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete post"})
		}
		return RenderTempl(c, posts.Error("Failed to delete post"))
	}

	if isJSON {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Post deleted"})
	}
	return c.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*Posts)(nil)
var _ CompleteResourceController = (*Posts)(nil)
