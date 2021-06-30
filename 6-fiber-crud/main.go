package main

import (
	"fibercrud/controllers"
	"fibercrud/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// *******************************************
func setupRoutes(app fiber.App) {

	// give response when at /
	app.Get("/", controllers.DefaultEndpoint)

	api := app.Group("/api")

	// give response when at /api
	api.Get("", controllers.ApiEndpoint)

	// connect getArticles over /api/articles
	routes.ArticlesRoutes(api.Group("/articles"))
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(*app)
	app.Listen(":8080")

}
