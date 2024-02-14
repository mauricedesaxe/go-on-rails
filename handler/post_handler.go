package handler

import (
	"goblog/model"
	"goblog/view/posts"
	"net/http"
	"strconv"

	"github.com/anthdm/slick"
)

type Posts struct{} // Implements RouteRegistrar and CoreHandler

func (h *Posts) RegisterRoutes(app *slick.Slick) {
	app.Get("/posts", h.index)
	app.Get("/posts/:id", h.show)
	app.Get("/posts/new", h.new)
	app.Post("/posts", h.create)
	app.Get("/posts/:id/edit", h.edit)
	app.Put("/posts/:id", h.update)
	app.Delete("/posts/:id", h.delete)
}

// GET /posts - index - List all posts
func (h *Posts) index(c *slick.Context) error {
	var rows []model.Post
	tx := model.DB.Find(&rows)
	if tx.Error != nil {
		return c.Render(posts.Error("No posts found"))
	}

	return c.Render(posts.Index(rows))
}

// GET /posts/:id - show - Show a single post
func (h *Posts) show(c *slick.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Render(posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Render(posts.Error("Post not found"))
	}

	return c.Render(posts.Show(res))
}

// GET /posts/new - new - Show the form to create a new post
func (h *Posts) new(c *slick.Context) error {
	return c.Render(posts.New())
}

// POST /posts - create - Create a new post
func (h *Posts) create(c *slick.Context) error {
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
		return c.Render(posts.Error("Failed to create post"))
	}

	return c.Redirect("/posts/"+strconv.Itoa(int(post.ID)), http.StatusFound)
}

// GET /posts/:id/edit - edit - Show the form to edit a post
func (h *Posts) edit(c *slick.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Render(posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Render(posts.Error("Post not found"))
	}

	return c.Render(posts.Edit(res))
}

// PUT /posts/:id - update - Update a post
func (h *Posts) update(c *slick.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Render(posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Render(posts.Error("Post not found"))
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	author := c.FormValue("author")
	res.Title = title
	res.Content = content
	res.Author = author
	tx = model.DB.Save(&res)
	if tx.Error != nil {
		return c.Render(posts.Error("Failed to update post"))
	}

	return c.Redirect("/posts/"+strconv.Itoa(int(res.ID)), http.StatusFound)
}

// DELETE /posts/:id - delete - Delete a post
func (h *Posts) delete(c *slick.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Render(posts.Error("Invalid ID format"))
	}

	var res model.Post
	tx := model.DB.First(&res, id)
	if tx.Error != nil {
		return c.Render(posts.Error("Post not found"))
	}

	tx = model.DB.Delete(&res)
	if tx.Error != nil {
		return c.Render(posts.Error("Failed to delete post"))
	}

	return c.Redirect("/posts", http.StatusFound)
}

// Ensure that Posts implements RouteRegistrar and CoreHandler
var _ RouteRegistrar = (*Posts)(nil)
var _ CompleteResourceController = (*Posts)(nil)
