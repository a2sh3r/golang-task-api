package server

import (
	"github.com/a2sh3r/golang-task-api.git/internal/middleware"
	"github.com/a2sh3r/golang-task-api.git/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	taskService service.TaskServiceInterface
}

func NewHandler(taskService service.TaskServiceInterface) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func NewRouter(handler *Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(middleware.NewLoggingMiddleware())
	app.Use(middleware.NewGzipMiddleware())

	app.Route("/tasks", func(router fiber.Router) {
		router.Post("/", handler.CreateTask)
		router.Get("/", handler.GetAllTasks)
		router.Get("/:id", handler.GetTask)
		router.Put("/:id", handler.UpdateTask)
		router.Delete("/:id", handler.DeleteTask)
	})

	return app
}
