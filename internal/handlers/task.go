package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"

	"github.com/cabrn3t/todo-list-server/internal/models"
	"github.com/cabrn3t/todo-list-server/internal/repositories"
	"github.com/gofiber/fiber/v3"
)

type TaskHandler struct {
	taskRepository repositories.Repository
}

func NewTaskHandler(taskRepository repositories.Repository) *TaskHandler {
	return &TaskHandler{taskRepository: taskRepository}
}

func (t *TaskHandler) CreateTask(ctx fiber.Ctx) error {
	task := models.Task{}
	json.Unmarshal(ctx.Body(), &task)

	if task.Title == "" {
		err := errors.New("title is empty")
		slog.Any("error", err)
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	switch task.Status {
	case "new", "":
		task.Status = models.NEW
	case "in_progress":
		task.Status = models.IN_PROGRESS
	case "done":
		task.Status = models.DONE
	default:
		err := errors.New("invalid status")
		slog.Any("error", err)
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	err := t.taskRepository.CreateTask(ctx.Context(), &task)
	if err != nil {
		slog.Error("error while inserting task", slog.Any("err", err))
		return ctx.Status(fiber.ErrInternalServerError.Code).SendString(err.Error())
	}
	return nil
}

func (t *TaskHandler) GetTasks(ctx fiber.Ctx) error {
	tasks, err := t.taskRepository.GetTasks(ctx.Context())
	if err != nil {
		slog.Error("error while receiving task slice", slog.Any("err", err))
		return ctx.Status(fiber.ErrInternalServerError.Code).SendString(err.Error())
	}
	err = ctx.JSON(tasks)
	if err != nil {
		slog.Error("error while forming JSON", slog.Any("err", err))
		return ctx.Status(fiber.ErrInternalServerError.Code).SendString(err.Error())
	}
	return nil
}

func (t *TaskHandler) UpdateTask(ctx fiber.Ctx) error {
	task := models.Task{}
	json.Unmarshal(ctx.Body(), &task)
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		err := errors.New("invalid status")
		slog.Any("error", err)
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	task.ID = int64(id)
	switch task.Status {
	case "":
		break
	case "new":
		task.Status = models.NEW
	case "in_progress":
		task.Status = models.IN_PROGRESS
	case "done":
		task.Status = models.DONE
	default:
		err := errors.New("invalid status")
		slog.Any("error", err)
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	err = t.taskRepository.UpdateTask(ctx.Context(), &task)
	if err != nil {
		slog.Error("error while updating task", slog.Any("err", err))
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return nil
}

func (t *TaskHandler) DeleteTask(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		err := errors.New("invalid id")
		slog.Any("error", err)
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	err = t.taskRepository.DeleteTask(ctx.Context(), int64(id))
	if err != nil {
		slog.Error("error while deleting task", slog.Any("err", err))
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return nil
}
