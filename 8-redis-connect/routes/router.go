package routes

import (
	"redisconnect/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteAPI(routeApi fiber.Router) {
	routeApi.Get("", controllers.GetArticles)
	routeApi.Post("", controllers.AddArticle)
	routeApi.Get("/:id", controllers.GetArticleById)
}

func CheckDB(apiRedis fiber.Router) {
	apiRedis.Get("/ping", controllers.PingRedis)
	apiRedis.Get("/:id", controllers.GetArticleByIdRedis)
}
