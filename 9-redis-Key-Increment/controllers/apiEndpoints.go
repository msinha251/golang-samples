package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var contextTimeout = 120 * time.Second

type SiteIdIncreased struct {
	Increased bool
}

// MAP REDIS KEY TO MONGODB via GOLANG DICTIONARY (MAP)
func UpdateRedisKeysInMongo() map[string]string {
	var siteIdDict = make(map[string]string)
	// ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	rds := RedisConnect()
	fmt.Println(rds.Ping(ctx).Val())

	client := MongoConnect(ctx)
	increasedBool := client.Database("go-crud").Collection("siteidincreased").FindOne(ctx, bson.D{{}})
	var siteidincreased SiteIdIncreased
	increasedBool.Decode(&siteidincreased)
	fmt.Println(siteidincreased)

	if rds.Ping(ctx).Val() != "PONG" {
		fmt.Println("Redis is down")
		return siteIdDict
	} else if rds.Ping(ctx).Val() == "PONG" && siteidincreased.Increased == false {
		rdsKeys := rds.Keys(ctx, "*")
		fmt.Println("rds keys are ", rdsKeys.Val())
		for _, key := range rdsKeys.Val() {
			siteIdDict[string(key)] = rds.Get(ctx, string(key)).Val()
		}

		// Update siteIdDict in mongodb
		client := MongoConnect(ctx)
		documents, err := client.Database("go-crud").Collection("siteid").CountDocuments(ctx, bson.D{})
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
		fmt.Println(documents)
		fmt.Println(siteIdDict)
		if documents == 0 {
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
		} else if len(siteIdDict) != 0 {
			// client.Database("go-crud").Collection("siteid").FindOneAndReplace(ctx, bson.D{}, siteIdDict)
			// replacing all documents in siteid collection
			client := MongoConnect(ctx)
			fmt.Println("Updated siteIdDict is ", siteIdDict)
			// deleting old siteid mongo document
			client.Database("go-crud").Collection("siteid").FindOneAndDelete(ctx, bson.D{})
			// Inserting new mongo document with Increatmented SiteID
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
			fmt.Println("MONGODB IS UPDATED with increased siteid under siteid collection")
			defer func() {
				if err := client.Disconnect(ctx); err != nil {
					panic(err)
				}
			}()
		} else {
			fmt.Println("Either siteid collection or redisKeys are empty")
			return siteIdDict
		}
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// defer func() {
	// 	if err := rds.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()
	fmt.Println(siteIdDict)
	return siteIdDict
}

// INCREASE KEYS WITH 500 IF REDIS IS DOWN
func IncrementKey() string {
	var siteIdDict = make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
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
		fmt.Println("Updated siteIdDict is ", siteIdDict)
		if documents == 0 {
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
		} else if len(siteIdDict) != 0 {
			// replacing all documents in siteid collection
			client := MongoConnect(ctx)
			fmt.Println("Updated siteIdDict is ", siteIdDict)
			// deleting old siteid mongo document
			client.Database("go-crud").Collection("siteid").FindOneAndDelete(ctx, bson.D{})
			// Inserting new mongo document with Increatmented SiteID
			client.Database("go-crud").Collection("siteid").InsertOne(ctx, siteIdDict)
			fmt.Println("MONGODB IS UPDATED with increased siteid under siteid collection")
			return "MONGODB IS UPDATED with increased siteid under siteid collection"
		} else {
			fmt.Println("Either siteid collection or redisKeys are empty, NO UPDATE IN MONGODB")
		}
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Redis is UP OR SiteId already Increased.")
	return "Redis is UP OR SiteId already Increased."
}

// UPDATE REDIS KEYS FROM MONGO WITH INCREATMENT KEYS
func UpdateRedisKeyFromMongo() string {
	var siteIdDict = make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Grabbing the sitesids mongo
	client := MongoConnect(ctx)
	mongoSiteIds, err := client.Database("go-crud").Collection("siteid").Find(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
	defer mongoSiteIds.Close(ctx)
	// Update REDIS KEYS
	for mongoSiteIds.Next(ctx) {
		mongoSiteIds.Decode(&siteIdDict)
		fmt.Println("Mongo SiteIDs", siteIdDict)
		break
	}

	// Grabbing the sitesidincreased Bool value from mongo

	increasedBool := client.Database("go-crud").Collection("siteidincreased").FindOne(ctx, bson.D{{}})
	var siteidincreased SiteIdIncreased
	increasedBool.Decode(&siteidincreased)
	fmt.Println(siteidincreased)
	rds := RedisConnect()

	// Loop though siteIDs
	if rds.Ping(ctx).Val() == "PONG" && siteidincreased.Increased == true {
		for country := range siteIdDict {
			if country != "_id" {
				rds.Set(ctx, country, siteIdDict[country], 0)
				fmt.Println("Redis Key " + country + " is Updated / Reinitialized with " + siteIdDict[country])
			}
		}
		siteidincreased.Increased = false
		client.Database("go-crud").Collection("siteidincreased").FindOneAndReplace(ctx, bson.M{"increased": true}, bson.M{"increased": false})
		fmt.Println("Redis Key has been Updated / Reinitialized with 500 increment, Also siteidincreased has been turned to FALSE")
		return "Redis Key has been Updated / Reinitialized with 500 increment, Also siteidincreased has been turned to FALSE"
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Redis is UP or SiteID already Increased")
	return "Redis is UP or SiteID already Increased"

}
