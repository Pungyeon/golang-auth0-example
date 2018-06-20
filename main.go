package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Pungyeon/golang-auth0-example/handlers"
)

func main() {
	usageFlag := flag.Bool("h", false, "[bool] show usage")
	flag.Parse()

	if *usageFlag {
		usage()
		os.Exit(1)
	}

	r := gin.Default()

	// r := mux.NewRouter()
	r.GET("/", handlers.IndexHandler)
	r.GET("/callback", handlers.CallbackHandler)
	r.GET("/user", handlers.UserHandler)
	r.GET("/login", handlers.LoginHandler)

	err := r.Run(":3000")
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
