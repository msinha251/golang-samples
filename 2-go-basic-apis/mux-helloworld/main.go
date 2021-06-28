package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func hellomux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Mux")
}

func handlers() {
	Router := mux.NewRouter()
	Router.HandleFunc("/hellomux", hellomux)
	http.ListenAndServe(":8080", Router)
}

func main() {
	handlers()
}
