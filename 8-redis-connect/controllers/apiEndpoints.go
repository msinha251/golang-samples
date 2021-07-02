package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Article struct {
	Id      int    `json:id`
	Title   string `json:title`
	Content string `json:content`
}

// GET ALL ARTICLES
func GetArticles(c *fiber.Ctx) error {
	var allArticles []Article

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := MongoConnect(ctx)
	articles, err := client.Database("go-crud").Collection("fiber-crud").Find(ctx, bson.D{})
	if err != nil {
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot connect to mongoDB",
		})
	}
	for articles.Next(ctx) {
		var article Article
		err := articles.Decode(&article)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"message": "Cannot decode articles",
			})
		} else {
			allArticles = append(allArticles, article)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"articles": allArticles,
		},
	})
}

// ALL NEW ARTICLE
func AddArticle(c *fiber.Ctx) error {
	var article Article
	c.BodyParser(&article)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := MongoConnect(ctx)
	articleCounts, err := client.Database("go-crud").Collection("fiber-crud").CountDocuments(ctx, bson.D{})
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot fetch total number of Documents.",
		})
	}
	id := int(articleCounts) + 1
	article.Id = id
	client.Database("go-crud").Collection("fiber-crud").InsertOne(ctx, article)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Article " + strconv.Itoa(id) + " ADDED",
	})

}

// GET ARTICLE BY ID
func GetArticleById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := MongoConnect(ctx)
	foundArticle := client.Database("go-crud").Collection("fiber-crud").FindOne(ctx, bson.M{"id": id})
	var article Article
	foundArticle.Decode(&article)
	if article.Id == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Article " + paramId + " NOT FOUND",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"body": fiber.Map{
			"article": article,
		},
	})
}

// ************** REDIS ****************
func PingRedis(c *fiber.Ctx) error {
	ctx := context.Background()
	rdb := RedisConnect()
	res := rdb.Ping(ctx)
	if res.Val() == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Unable to connect with Redis, please check Redis is running or not.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": res.Val(),
	})
}

func GetArticleByIdRedis(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse Id.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rds := RedisConnect()
	articleID := rds.Keys(ctx, paramId)
	fmt.Println(articleID)
	fmt.Println(paramId)
	if len(articleID.Val()) == 0 {
		// connect to mongo and fetch article
		client := MongoConnect(ctx)
		foundArticle := client.Database("go-crud").Collection("fiber-crud").FindOne(ctx, bson.M{"id": id})
		var article Article
		foundArticle.Decode(&article)
		jsonarticle, err := json.Marshal(article)

		fmt.Println()
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": false,
				"message": "Cannot Marshal article (struct) to json",
			})
		}
		// push the article to redis cache
		rds.Set(ctx, paramId, string(jsonarticle), 0)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"article": article,
		})
	}

	rdsArticle := rds.Get(ctx, paramId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"article": rdsArticle.Val(),
	})
}
