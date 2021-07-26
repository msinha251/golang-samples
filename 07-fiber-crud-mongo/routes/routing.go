package routes

import (
	"fibercrudmongo/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(route fiber.Router) {
	route.Get("", controllers.GetArticles)
	route.Post("", controllers.AddArticles)
	route.Delete("/:id", controllers.DeleteArticleById)
	route.Get("/:id", controllers.GetArticleById)
	route.Put("/:id", controllers.UpdateArticleById)
}
