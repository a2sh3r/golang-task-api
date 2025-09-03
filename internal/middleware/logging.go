package middleware

import (
	"time"

	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		logger.Log.Info("HTTP request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int("size", len(c.Response().Body())),
			zap.Duration("duration", duration),
		)

		return err
	}
}
