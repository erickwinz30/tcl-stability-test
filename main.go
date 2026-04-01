package main

import (
	"stability-test-task-api/handlers"
	"stability-test-task-api/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New()

	// menambahkan rate limiter untuk mencegah abuse
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				types.NewErrorResponse("too many requests, please try again later", nil),
			)
		},
	}))

	app.Get("/tasks", handlers.GetTasks)
	app.Get("/tasks/:id", handlers.GetTask)
	app.Post("/tasks", handlers.CreateTask)
	app.Delete("/tasks/:id", handlers.DeleteTask)
	// endpoint baru untuk update
	app.Put("/tasks/:id", handlers.UpdateTask)

	app.Listen(":3000")
}
