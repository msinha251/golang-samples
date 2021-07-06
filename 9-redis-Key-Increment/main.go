package main

import (
	"fmt"
	"redisconnect/controllers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ******* Function for calling every 2 seconds *********

func runEverySecond() {
	for {
		time.Sleep(2 * time.Second)
		go controllers.UpdateRedisKeysInMongo()
		fmt.Println("Push redis key in mongo")
		go controllers.IncrementKey()
		fmt.Println("Incrementing keys")
		go controllers.UpdateRedisKeyFromMongo()
		fmt.Println("Update redis key from mongo")
		fmt.Println("****************************")
	}

}

// func startPolling1() {
// 	for {
// 		time.Sleep(2 * time.Second)
// 		go doSomething("from polling 1")
// 	}
// }

// ***********************************************

func main() {
	app := fiber.New()
	app.Use(logger.New())
	// setupRoutes(app)
	// go runEverySecond()

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println(controllers.UpdateRedisKeysInMongo())
			time.Sleep(1 * time.Second)

			fmt.Println(controllers.IncrementKey())
			time.Sleep(1 * time.Second)

			fmt.Println(controllers.UpdateRedisKeyFromMongo())
			time.Sleep(1 * time.Second)
			// fmt.Println("hello")
		}
	}()

	// wait for 10 seconds
	// time.Sleep(5 * time.Second)
	// ticker.Stop()
	// fmt.Println(controllers.UpdateRedisKeysInMongo())
	// fmt.Println(controllers.IncrementKey())
	// // fmt.Println(controllers.UpdateRedisKeyFromMongo())

	app.Listen(":8080")
}
