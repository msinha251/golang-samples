package routes

import (
	"fibercrud/controllers"

	"github.com/gofiber/fiber/v2"
)

func ArticlesRoutes(route fiber.Router) {
	route.Get("", controllers.GetArticles)
	route.Post("", controllers.AddArticle)
	route.Get("/:id", controllers.GetArticleById)
	route.Delete("/:id", controllers.DeleteArticleById)
	route.Put("/:id", controllers.UpdateArticleById)
}
