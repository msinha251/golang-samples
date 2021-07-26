package main

import (
	"fibercrudmongo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {

	api := app.Group("/api")

	routes.Router(api.Group("/articles"))
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	app.Listen(":8080")
}
