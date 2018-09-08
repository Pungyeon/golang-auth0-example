package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/auth0-community/auth0"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/Pungyeon/golang-auth0-example/db"
	"github.com/Pungyeon/golang-auth0-example/handlers"
)

/* USAGE:
 * Set the environment variables:
 * 	AUTH0_CLIEN_ID
 * 	AUTH0_DOMAIN: "https://pungy.eu.auth0.com/" (very importantly not omitting the last /)
 */

var (
	audience    string
	domain      string
	todoHandler *handlers.TodoHandler
)

func main() {
	backend := flag.String("db", "memory", "specify which backend to use for our todo ap: ('memory' | 'postgres')")
	config := flag.String("config", "config.json", "specify the relative filepath of the config to use for the postgres db connecdtion")
	flag.Parse()

	switch *backend {
	case "memory":
		todoHandler = handlers.NewTodoHandler(
			db.NewInMemoryDB(),
		)
	case "postgres":
		pqConfig, err := db.ReadPostgreConfig(*config)
		if err != nil {
			panic(err)
		}
		todoHandler = handlers.NewTodoHandler(
			db.NewPostgreTodoDB(pqConfig),
		)
	}

	setAuth0Variables()
	r := gin.Default()

	// This will ensure that the angular files are served correctly
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	authorized := r.Group("/")
	authorized.Use(authRequired())
	authorized.GET("/todo", todoHandler.GetTodoListHandler)
	authorized.POST("/todo", todoHandler.AddTodoHandler)
	authorized.DELETE("/todo/:id", todoHandler.DeleteTodoHandler)
	authorized.PUT("/todo", todoHandler.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_CLIENT_ID")
	domain = os.Getenv("AUTH0_DOMAIN")
}

// ValidateRequest will verify that a token received from an http request
// is valid and signy by authority
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
