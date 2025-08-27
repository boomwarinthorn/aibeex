package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create fiber app
	app := fiber.New()

	// GET endpoint that returns JSON "Hello World"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	// Start server on port 3000
	app.Listen(":3000")
}
