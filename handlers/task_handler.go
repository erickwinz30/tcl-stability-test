package handlers

import (
	"errors"
	"strconv"

	"stability-test-task-api/models"
	"stability-test-task-api/store"
	"stability-test-task-api/types"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func GetTasks(c *fiber.Ctx) error {
	tasks := store.GetAllTasks()
	return c.Status(fiber.StatusOK).JSON(types.NewSuccessResponse("tasks fetched", tasks))
}

func GetTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("invalid task id", err.Error()))
	}

	task := store.GetTaskByID(id)

	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(types.NewErrorResponse("task not found", nil))
	}

	return c.Status(fiber.StatusOK).JSON(types.NewSuccessResponse("task fetched", task))
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("cannot parse json body", err.Error()))
	}

	// validasi menggunakan go-playground/validator
	if err := validate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("validation failed", err.Error()))
	}

	// membuat logic generate id baru yang sebelumnya belum ada
	tasks := store.GetAllTasks()
	newID := 1

	if len(tasks) > 0 {
		maxID := tasks[0].ID
		for _, t := range tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		newID = maxID + 1
	}

	task.ID = newID

	store.AddTask(task)

	return c.Status(fiber.StatusCreated).JSON(types.NewSuccessResponse("task created", task))
}

func DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("invalid task id", err.Error()))
	}

	task := store.GetTaskByID(id)
	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(types.NewErrorResponse("task not found", nil))
	}

	store.DeleteTask(id)

	return c.Status(fiber.StatusOK).JSON(types.NewSuccessResponse("task deleted", nil))
}

func UpdateTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("invalid task id", err.Error()))
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("cannot parse json body", err.Error()))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.NewErrorResponse("validation failed", err.Error()))
	}

	updatedTask, err := store.UpdateTask(id, req.Title, *req.Done)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.NewErrorResponse("task not found", nil))
		}
		if errors.Is(err, store.ErrNoChanges) {
			return c.Status(fiber.StatusConflict).JSON(types.NewErrorResponse("no changes detected", nil))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(types.NewErrorResponse("failed to update task", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(types.NewSuccessResponse("task updated", updatedTask))
}
