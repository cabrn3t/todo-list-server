package routes

import (
	"github.com/cabrn3t/todo-list-server/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

type TaskRoutes struct {
	taskHandler handlers.TaskHandler
}

func NewTaskRoutes(taskHandler handlers.TaskHandler) *TaskRoutes {
	return &TaskRoutes{taskHandler: taskHandler}
}

func (r *TaskRoutes) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/tasks", func(ctx fiber.Ctx) error {
		return r.taskHandler.CreateTask(ctx)
	})
	api.Get("/tasks", func(ctx fiber.Ctx) error {
		return r.taskHandler.GetTasks(ctx)
	})
	api.Put("/tasks/:id", func(ctx fiber.Ctx) error {
		return r.taskHandler.UpdateTask(ctx)
	})
	api.Delete("/tasks/:id", func(ctx fiber.Ctx) error {
		return r.taskHandler.DeleteTask(ctx)
	})
}
