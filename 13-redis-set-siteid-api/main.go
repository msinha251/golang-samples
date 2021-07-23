package main

import (
	"fmt"
	"log"
	"os"
	"redissetsiteids/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	routes.Router(api.Group("/redis"))
}

func main() {

	// environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	fmt.Println(os.Getenv("SITE_TITLE"))
	fmt.Println(os.Getenv("REDIS_HOST"))

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	app.Listen(":8080")
}
