package routes

import (
	"redissetsiteids/connections"

	"github.com/gofiber/fiber/v2"
)

func Router(route fiber.Router) {
	route.Post("", connections.SetSiteId)
}
