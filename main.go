package main

import (
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/gbrlsnchs/jwt"
	"github.com/gin-gonic/gin"

	"github.com/Pungyeon/golang-auth0-example/handlers"
)

// ValidateRequest will verify that a token received from an http request
// is valid and signy by authority
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		jot, err := jwt.FromRequest(c.Request)
		if err != nil {
			terminateWithError(http.StatusUnauthorized, "cannot read token from request", c)
			return
		}

		if jot.Validate(jwt.AlgorithmValidator("RS256")) != nil {
			terminateWithError(http.StatusUnauthorized, "could not validate signing algorithm", c)
			return
		}

		if jot.Validate(jwt.IssuerValidator("https://pungy.eu.auth0.com/")) != nil {
			terminateWithError(http.StatusUnauthorized, "could not validate issuer", c)
			return
		}

		if jot.Validate(jwt.ExpirationTimeValidator(time.Now())) != nil {
			terminateWithError(http.StatusUnauthorized, "token has expired", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func main() {
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
	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddTodoHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
