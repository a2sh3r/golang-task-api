package server

import (
	"github.com/a2sh3r/golang-task-api.git/internal/logger"
	"github.com/a2sh3r/golang-task-api.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	task, err := h.taskService.CreateTask(c.Context(), req.Title, req.Description)
	if err != nil {
		logger.Log.Error("failed to create task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *Handler) GetAllTasks(c *fiber.Ctx) error {
	tasks, err := h.taskService.GetAllTasks(c.Context())
	if err != nil {
		logger.Log.Error("failed to get tasks", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tasks)
}

func (h *Handler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	task, err := h.taskService.GetTask(c.Context(), id)
	if err != nil {
		logger.Log.Error("failed to get task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.JSON(task)
}

func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	task, err := h.taskService.UpdateTask(c.Context(), id, req)
	if err != nil {
		logger.Log.Error("failed to update task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(task)
}

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	err := h.taskService.DeleteTask(c.Context(), id)
	if err != nil {
		logger.Log.Error("failed to delete task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}
