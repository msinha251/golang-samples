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

// GET ARTICLE BY ID From Cache
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

// MAP REDIS KEY TO MONGODB via GOLANG DICTIONARY (MAP)
func UpdateRedisKeysInMongo() map[string]string {
	siteIdDict := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rds := RedisConnect()
	fmt.Println(rds.Ping(ctx).Val())
	if rds.Ping(ctx).Val() != "PONG" {
		fmt.Println("Redis is down")
		return siteIdDict
	}

	rdsKeys := rds.Keys(ctx, "*")
	// fmt.Println(rdsKeys.Val())
	for _, key := range rdsKeys.Val() {
		siteIdDict[string(key)] = rds.Get(ctx, string(key)).Val()
	}

	// Update siteIdDict in mongodb
	client := MongoConnect(ctx)
	documents, err := client.Database("go-crud").Collection("siteid").CountDocuments(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println(documents)
	fmt.Println(siteIdDict)
	if documents == 0 {
		client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
	} else if len(siteIdDict) != 0 {
		client.Database("go-crud").Collection("siteid").FindOneAndReplace(ctx, bson.D{}, siteIdDict)
	} else {
		fmt.Println("Either siteid collection or redisKeys are empty")
		return siteIdDict
	}
	fmt.Println()
	return siteIdDict
}

// INCREASE KEYS WITH 500 IF REDIS IS DOWN
func IncrementKey() string {
	siteIdDict := make(map[string]string)
	type SiteIdIncreased struct {
		Increased bool
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Grabbing the sitesidincreased Bool value from mongo
	client := MongoConnect(ctx)
	increasedBool := client.Database("go-crud").Collection("siteidincreased").FindOne(ctx, bson.D{{}})
	var siteidincreased SiteIdIncreased
	increasedBool.Decode(&siteidincreased)
	fmt.Println(siteidincreased)

	rds := RedisConnect()

	// Increase siteID by 500 if Redis is down and not already increased
	if rds.Ping(ctx).Val() != "PONG" && siteidincreased.Increased == false {
		siteid, err := client.Database("go-crud").Collection("siteid").Find(ctx, bson.D{{}})
		if err != nil {
			panic(err)
		}
		// SiteID increment
		for siteid.Next(ctx) {
			siteid.Decode(&siteIdDict)
			fmt.Println(siteIdDict)
			for country := range siteIdDict {
				if country != "_id" {
					var id int
					id, err = strconv.Atoi(string(siteIdDict[country]))
					if err != nil {
						panic(err)
					}
					siteIdDict[country] = strconv.Itoa(id + 500)
				}
			}
		}

		// Changing SiteIdIncreased to TRUE, hence no more futher increment is required.
		siteidincreased.Increased = true
		client.Database("go-crud").Collection("siteidincreased").FindOneAndReplace(ctx, bson.M{"increased": false}, bson.M{"increased": true})

		// Update increased SiteIdDict in mongodb
		documents, err := client.Database("go-crud").Collection("siteid").CountDocuments(ctx, bson.D{})
		if err != nil {
			panic(err)
		}
		fmt.Println("documents are ", documents)
		fmt.Println("Updated siteIdDict is ", siteIdDict)
		if documents == 0 {
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
		} else if len(siteIdDict) != 0 {
			// replacing all documents in siteid collection
			client := MongoConnect(ctx)
			fmt.Println("Updated siteIdDict is ", siteIdDict)
			// testing
			deleted := client.Database("go-crud").Collection("siteid").FindOneAndDelete(ctx, bson.D{})
			fmt.Println(deleted)
			documents2, err := client.Database("go-crud").Collection("siteid").CountDocuments(ctx, bson.D{})
			if err != nil {
				panic(err)
			}
			fmt.Println(documents2)
			// testing
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
			return "MONGODB IS UPDATED with increased siteid under siteid collection"
		} else {
			fmt.Println("Either siteid collection or redisKeys are empty, NO UPDATE IN MONGODB")
		}
	}
	return "Redis is UP OR SiteId already Increased."

}
