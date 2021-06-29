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

func getDatabasesCount(w http.ResponseWriter, r *http.Request) {
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
	var newArticle Article
	reqBody, _ := ioutil.ReadAll(r.Body)

	// ********* Content and MongoConnect **********
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// *********************************************

	json.Unmarshal(reqBody, &newArticle)
	// article := Article{
	// 	Id:      "1",
	// 	Title:   "Title 1",
	// 	Content: "my articles content",
	// }

	// check if the article is already available or not
	document := client.Database("go-crud").Collection("go-crud").FindOne(ctx, bson.D{{Key: "id", Value: newArticle.Id}})
	var article Article
	document.Decode(&article)

	if article.Id != "" {
		fmt.Fprintf(w, "Article CANNOT be added, since ID is already available, please try with new ID.")
	} else {
		client.Database("go-crud").Collection("go-crud").InsertOne(ctx, newArticle)
		fmt.Fprintf(w, "Article added.")
	}
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	var results []Article
	// ********* Content and MongoConnect **********
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// *********************************************

	articles, err := client.Database("go-crud").Collection("go-crud").Find(ctx, bson.D{})

	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer articles.Close(ctx)
	} else {
		for articles.Next(ctx) {
			var result Article
			err := articles.Decode(&result)
			if err != nil {
				fmt.Println("curson.Next() ERROR: ", err)
			} else {
				// fmt.Println(results)
				// json.NewEncoder(w).Encode(result)
				results = append(results, result)
			}
		}
	}
	json.NewEncoder(w).Encode(results)
}

func getArticleById(w http.ResponseWriter, r *http.Request) {

	// ********* Content and MongoConnect **********
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// *****************************************
	vars := mux.Vars(r)
	id := vars["id"]

	document := client.Database("go-crud").Collection("go-crud").FindOne(ctx, bson.D{{Key: "id", Value: id}})

	var article Article
	document.Decode(&article)

	if article.Id == "" {
		fmt.Fprintf(w, "Article Not Found for ID "+id)
	} else {
		json.NewEncoder(w).Encode(article)
	}
}

func handlers() {
	router := mux.NewRouter()
	router.HandleFunc("/getdatabasescount", getDatabasesCount)
	router.HandleFunc("/addarticle", addArticle).Methods("POST")
	router.HandleFunc("/getarticles", getArticles)
	router.HandleFunc("/getarticlebyid/{id}", getArticleById)
	http.ListenAndServe(":8080", router)
}

// *****************************************

func main() {
	handlers()
}
