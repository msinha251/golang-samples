package main

import (
	"context"
	"fmt"
	"log"
	"mongoautoincreatment/handlers"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type employee struct {
	id   int    `json:id`
	name string `name:`
}

func mongoConnect(ctx context.Context) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	documents, err := client.Database("test").Collection("test").CountDocuments(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strconv.Itoa(int(documents)))

	// insert counter collection
	client.Database("test").Collection("test").InsertOne(ctx, bson.D{
		{Key: "_id", Value: "userid"},
		{Key: "seq", Value: 0},
	})

}

func getNextSequence(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ret := client.Database("test").Collection("test").FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: name}}, bson.D{{Key: "$inc", Value: bson.M{"seq": 1}}})
	fmt.Printf("ret.Decode(&test): %v\n", ret)
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	handlers.Handlers()
	// mongoConnect(ctx)
	// getNextSequence("Mahesh")
}
