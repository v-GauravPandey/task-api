package http

import (
	"github.com/gofiber/fiber/v2"
)

// NewApp creates and configures a new Fiber application.
func NewApp(handler *TaskHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		},
	})

	// Register API routes directly (routes use /tasks path)
	handler.RegisterRoutes(app)

	return app
}
