package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Pungyeon/golang-auth0-example/handlers"
	"github.com/gorilla/mux"
)

func main() {
	usageFlag := flag.Bool("h", false, "[bool] show usage")
	flag.Parse()

	if *usageFlag {
		usage()
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/callback", handlers.CallbackHandler)
	r.HandleFunc("/user", handlers.UserHandler)
	r.HandleFunc("/login", handlers.LoginHandler)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Printf(`
		The following environment variables are required:
		AUTH0_COOKIE_SECRET
		AUTH0_DOMAIN
		AUTH0_CLIENT_ID
		AUTH0_CLIENT_SECRET
		AUTH0_CALLBACK_URL
	`)
}
