package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Pungyeon/golang-auth0-example/db"
	"github.com/Pungyeon/golang-auth0-example/handlers"
)

/* USAGE:
 * Set the environment variables:
 * 	AUTH0_CLIEN_ID
 * 	AUTH0_DOMAIN: "https://pungy.eu.auth0.com/" (very importantly not omitting the last /)
 */

var (
	todoHandler *handlers.TodoHandler
)

func main() {
	auth0Audience := flag.String("audience", "", "specify client id for connecting to auth0")
	auth0Domain := flag.String("domain", "", "specify the domain for connecting to auth0")
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

	authHandler := handlers.NewAuthHandler(
		DetermineAuth0Variables(*auth0Audience, *auth0Domain),
	)

	authorized := r.Group("/")
	authorized.Use(authHandler.Required())

	authorized.GET("/todo", todoHandler.GetTodoListHandler)
	authorized.POST("/todo", todoHandler.AddTodoHandler)
	authorized.DELETE("/todo/:id", todoHandler.DeleteTodoHandler)
	authorized.PUT("/todo", todoHandler.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

// DetermineAuth0Variables will set the domain and audience values, based on CLI input
// and if this input is empty, will try to retrieve these values from environment variables
func DetermineAuth0Variables(audience string, domain string) (string, string) {
	if audience == "" {
		log.Println("Audience not detect from CLI parameters, attempting to retrieve from ENV variables")
		audience = os.Getenv("AUTH0_CLIENT_ID")
	}
	if domain == "" {
		log.Println("Domain not detect from CLI parameters, attempting to retrieve from ENV variables")
		domain = os.Getenv("AUTH0_DOMAIN")
	}
	return audience, domain
}
