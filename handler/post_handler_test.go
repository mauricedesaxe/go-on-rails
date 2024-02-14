package handler

import (
	"bytes"
	"goblog/model"
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
	model.DB = model.InitDB("store.test.db")
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

	// create a post with JSON body
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
	// Assuming the API returns the created post's details
	assert.JSONEq(t, postData, string(bodyBytes))

	// try reading all posts, assert that one exists
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
	assert.JSONEq(t, `[{"id":1,"title":"Test Post","content":"This is a test post","author":"Tester"}]`, string(bodyBytes))

	// try reading a single post, assert that it exists
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
	assert.JSONEq(t, `{"id":1,"title":"Test Post","content":"This is a test post","author":"Tester"}`, string(bodyBytes))

	// update a post with JSON body
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
	assert.JSONEq(t, updateData, string(bodyBytes))

	// delete a post
	req = httptest.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close() // No need to read body for DELETE as we expect no content
}
