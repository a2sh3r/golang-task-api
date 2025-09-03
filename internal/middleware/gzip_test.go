package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewGzipMiddleware(t *testing.T) {
	app := fiber.New()
	middleware := NewGzipMiddleware()

	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test response with longer content to ensure gzip compression is applied")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestNewGzipMiddleware_NoGzipHeader(t *testing.T) {
	app := fiber.New()
	middleware := NewGzipMiddleware()

	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test response")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("Content-Encoding"))
}

func TestNewGzipMiddleware_WithGzipBody(t *testing.T) {
	app := fiber.New()
	middleware := NewGzipMiddleware()

	app.Use(middleware)
	app.Post("/test", func(c *fiber.Ctx) error {
		body := c.Body()
		return c.SendString("received: " + string(body))
	})

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	gw.Write([]byte("compressed data"))
	gw.Close()

	req := httptest.NewRequest("POST", "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
