package connections

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetSiteId(c *fiber.Ctx) error {

	// Get the site-id json from the request
	var body = c.Body()
	var response = make(map[string]string)
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	// Connect to Redis
	client, err := RedisConnect()
	if err != nil {
		return err
	}
	defer client.Close()

	// Setting the site id in Redis and preparing the response
	for k, v := range data {
		if client.Ping().Val() == "PONG" {
			o := strings.Split(client.Set(k, v.(string), 0).String(), ": ")
			response[o[0]] = o[1]
		} else {
			response["set "+k] = "ERROR: Redis is not running OR unable to connect"
		}
	}

	// Send the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})

}
