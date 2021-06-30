package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Article struct {
	Id      int    `json:id`
	Title   string `json:title`
	Content string `json:content`
}

var Articles = []Article{
	{
		Id:      1,
		Title:   "Walk the dog 🦮",
		Content: "This is content",
	},
}

// ************** API Endpoints **************

// Default
func DefaultEndpoint(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "You are at Default Endpoint.",
	})
}

// /api
func ApiEndpoint(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "You are at /api endpoint.",
	})
}

// GetArticles
func GetArticles(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"articles": Articles,
		},
	})
}

// Add article
func AddArticle(c *fiber.Ctx) error {
	type Request struct {
		Title   string `json:title`
		Content string `json:"content"`
	}
	var body Request
	err := c.BodyParser(&body)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}
	fmt.Println(body)
	article := &Article{
		Id:      len(Articles) + 1,
		Title:   body.Title,
		Content: body.Content,
	}
	Articles = append(Articles, *article)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"body": fiber.Map{
			"articles": Articles,
		},
	})
}

// Get Article by ID:
func GetArticleById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	for _, article := range Articles {
		if article.Id == id {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"data": fiber.Map{
					"article": article,
				},
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": false,
		"message": "Article not found.",
	})
}

// Delete Article by ID:
func DeleteArticleById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"message": "Article " + paramId + " DELETED",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": false,
		"message": "Article NOT found",
	})

}

// Update Article:
func UpdateArticleById(c *fiber.Ctx) error {
	type Request struct {
		Title   string `json:title`
		Content string `json:"content"`
	}
	var body Request
	paramId := c.Params("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	for _, article := range Articles {
		if article.Id == id {
			fmt.Println(article)
			article.Title = body.Title
			article.Content = body.Content
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"message": "Article UPDATED",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": false,
		"message": "Article Not Found",
	})
}
