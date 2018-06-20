package main

import (
	"net/http"

	"github.com/Pungyeon/go-test/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/user", handlers.UserHandler)
	r.HandleFunc("/login", handlers.LoginHandler)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
