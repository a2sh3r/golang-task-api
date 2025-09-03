package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggingMiddleware(t *testing.T) {
	app := fiber.New()
	middleware := NewLoggingMiddleware()

	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestNewLoggingMiddleware_WithError(t *testing.T) {
	app := fiber.New()
	middleware := NewLoggingMiddleware()

	app.Use(middleware)
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "test error")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
