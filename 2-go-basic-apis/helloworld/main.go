package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func handeler() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handeler()
}
