package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func helloworld(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Hello World",
	})
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", helloworld)

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
