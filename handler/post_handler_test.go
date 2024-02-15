package handler

import (
	"bytes"
	"encoding/json"
	models "goblog/models"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	// Setup
	models.DB = models.InitDB("store.test.db")
	defer func() {
		os.Remove("store.test.db")
	}()

	app := fiber.New()
	h := Posts{}
	h.RegisterRoutes(app)

	// try reading all posts, assert that none exist
	req := httptest.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	assert.JSONEq(t, "[]", string(bodyBytes))

	// create a post with JSON body, assert that it was created
	postData := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	req = httptest.NewRequest("POST", "/posts", strings.NewReader(postData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	post := models.Post{}
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post", post.Content)
	assert.Equal(t, "Tester", post.Author)

	// try reading all posts, assert that the created post exists
	req = httptest.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	posts := []models.Post{}
	err = json.Unmarshal(bodyBytes, &posts)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, posts, 1)
	assert.Equal(t, "Test Post", posts[0].Title)
	assert.Equal(t, "This is a test post", posts[0].Content)
	assert.Equal(t, "Tester", posts[0].Author)

	// try reading a single post, assert that the created post exists
	req = httptest.NewRequest("GET", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	post = models.Post{}
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post", post.Content)
	assert.Equal(t, "Tester", post.Author)

	// update a post with JSON body, assert that it was updated
	updateData := `{"title":"Updated Test Post","content":"This is an updated test post","author":"Tester"}`
	req = httptest.NewRequest("PUT", "/posts/1", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	post = models.Post{}
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Updated Test Post", post.Title)
	assert.Equal(t, "This is an updated test post", post.Content)
	assert.Equal(t, "Tester", post.Author)

	// delete a post
	req = httptest.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// try reading all posts, assert that none exist
	req = httptest.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	posts = []models.Post{}
	err = json.Unmarshal(bodyBytes, &posts)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, posts, 0)
}
