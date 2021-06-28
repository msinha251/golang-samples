package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	Id      string `bson:id`
	Title   string `bson:title`
	Content string `bson:content`
}

// ************ MONGO CONNECTION **************

func mongoConnect(ctx context.Context) (*mongo.Client, error) {
	var mongoConnectionError error

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		mongoConnectionError = err
		log.Fatal(err)
	}
	return client, mongoConnectionError
}

// *****************************************

// ************ API ENDPOINTS **************

func getDatabases(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	databases, _ := client.ListDatabaseNames(ctx, bson.D{})
	fmt.Fprintf(w, strconv.Itoa(len(databases)))
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	var article Article
	reqBody, _ := ioutil.ReadAll(r.Body)

	// ********* mongo **********
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// *************************

	json.Unmarshal(reqBody, &article)
	// article := Article{
	// 	Id:      "1",
	// 	Title:   "Title 1",
	// 	Content: "my articles content",
	// }

	client.Database("go-crud").Collection("go-crud").InsertOne(ctx, article)
	fmt.Fprintf(w, "Article added.")
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongoConnect(ctx)
	articles, err := client.Database("go-crud").Collection("go-crud").Find(ctx, bson.D{})
	// var allArticles []Article
	// fmt.Println(articles.All(ctx, &allArticles))

	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer articles.Close(ctx)
	} else {
		for articles.Next(ctx) {
			var results bson.D
			err := articles.Decode(&results)
			if err != nil {
				fmt.Println("curson.Next() ERROR: ", err)
			} else {
				// fmt.Println(results)
				json.NewEncoder(w).Encode(results)
			}
		}
	}

	// results, _ := json.Marshal(articles)
	// fmt.Fprintf(w, string(results))

}

func handlers() {
	router := mux.NewRouter()
	router.HandleFunc("/getdatabases", getDatabases)
	router.HandleFunc("/addarticle", addArticle).Methods("POST")
	router.HandleFunc("/getarticles", getArticles)
	http.ListenAndServe(":8080", router)
}

// *****************************************

func main() {
	handlers()
}
