package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

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

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	// r := mux.NewRouter()
	r.GET("/auth/", handlers.IndexHandler)
	r.GET("/auth/callback", handlers.CallbackHandler)
	r.GET("/auth/user", handlers.UserHandler)
	r.GET("/auth/login", handlers.LoginHandler)
	r.GET("/auth/check-auth", handlers.CheckAuthHandler)

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
