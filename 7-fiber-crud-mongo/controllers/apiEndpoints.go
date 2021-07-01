package controllers

import (
	"context"
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

// get articles from Mongo

func GetArticles(c *fiber.Ctx) error {
	var allArticles []Article
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	client := MongoConnect(ctx)
	articles, err := client.Database("go-crud").Collection("fiber-crud").Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	} else {
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
}

// Add Articles to Mongo:
func AddArticles(c *fiber.Ctx) error {
	var body Article
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot pard JSON",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := MongoConnect(ctx)

	articlesCount, err := client.Database("go-crud").Collection("fiber-crud").CountDocuments(ctx, bson.D{})
	body.Id = int(articlesCount) + 1
	client.Database("go-crud").Collection("fiber-crud").InsertOne(ctx, body)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Article Addedd",
	})
}

// Delete Article by ID from Mongo:
func DeleteArticleById(c *fiber.Ctx) error {
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
	foundArticle := client.Database("go-crud").Collection("fiber-crud").FindOne(ctx, bson.M{"id": int(id)})
	var article Article
	foundArticle.Decode(&article)
	if article.Id == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Article " + paramId + " not Found.",
		})
	} else {
		client.Database("go-crud").Collection("fiber-crud").DeleteOne(ctx, bson.M{"id": int(id)})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Article " + paramId + " DELETED",
		})
	}
}

// Get Article By ID:
func GetArticleById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannod Parse ID",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := MongoConnect(ctx)
	var article Article
	document := client.Database("go-crud").Collection("fiber-crud").FindOne(ctx, bson.M{"id": id})
	document.Decode(&article)

	if article.Id == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Article " + paramId + " not found.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"article": article,
		},
	})
}

// Update Article By Id:
func UpdateArticleById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID.",
		})
	}
	var updateArticle Article
	err = c.BodyParser(&updateArticle)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON.",
		})
	}
	updateArticle.Id = id
	fmt.Println(updateArticle)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := MongoConnect(ctx)
	document := client.Database("go-crud").Collection("fiber-crud").FindOne(ctx, bson.M{"id": id})
	var article Article
	document.Decode(&article)
	fmt.Println(article)
	if article.Id == 0 {
		client.Database("go-crud").Collection("fiber-crud").InsertOne(ctx, updateArticle)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Article " + paramId + " not found in DB, but added now.",
		})
	}
	// id_, err := primitive.ObjectIDFromHex("60dd7726625f388ec992ebab")
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot  assign id_",
		})
	}
	client.Database("go-crud").Collection("fiber-crud").FindOneAndUpdate(ctx, bson.M{"id": id}, bson.M{"$set": updateArticle})
	fmt.Println(article)
	fmt.Println(updateArticle)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Article " + paramId + " UPDATED.",
	})
}
