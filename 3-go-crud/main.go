package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:id,omitempty`
	Topic   string `json:topic,omitempty`
	Content string `json:content,omitempty`
}

var Articles []Article

// **************CRUD API's*****************

func getArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint HIT: getArticles")
	json.NewEncoder(w).Encode(Articles)
}

func getArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	articleFound := false
	for _, article := range Articles {
		if article.Id == id {
			articleFound = true
			json.NewEncoder(w).Encode(article)
			break
		}
	}
	if articleFound == true {
		fmt.Println("Article with id" + id + "found")
	} else {
		fmt.Fprintf(w, "Article not found in database, please check the ID")
	}

	fmt.Println("Endpoint HIT: getArticleById")

}

func addArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newArticle Article
	var articleAvailable = false
	json.Unmarshal(reqBody, &newArticle)

	for _, article := range Articles {
		if article.Id == newArticle.Id {
			articleAvailable = true
			break
		}
	}

	if articleAvailable == true {
		fmt.Fprintf(w, "Article already available")
	} else {
		Articles = append(Articles, newArticle)
		fmt.Fprintf(w, "Article ADDED")
	}
	// json.NewEncoder(w).Encode(newArticle)
	fmt.Println("Endpoint HIT: addArticle")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	articleFound := false

	for index, article := range Articles {
		if article.Id == id {
			articleFound = true
			Articles = append(Articles[:index], Articles[index+1:]...)
			break
		}
	}

	if articleFound == true {
		fmt.Fprintf(w, "Article with ID "+id+" deleted")
	} else {
		fmt.Fprintf(w, "Article with ID "+id+" not found.")
	}
}

func handlers() {
	router := mux.NewRouter()
	router.HandleFunc("/getarticles", getArticles)
	router.HandleFunc("/getarticle/{id}", getArticleById)
	router.HandleFunc("/addarticle", addArticle).Methods("POST")
	router.HandleFunc("/deletearticle/{id}", deleteArticle)

	http.ListenAndServe(":8080", router)
}

// *****************************************
func main() {
	Articles = []Article{
		Article{Id: "1", Topic: "Article 1", Content: "Article Content"},
		Article{Id: "2", Topic: "Article 2", Content: "Article Content"},
	}
	handlers()
}
