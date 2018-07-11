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
			c.JSON(http.StatusUnauthorized, "cannot read token from request")
			return
		}

		if jot.Validate(jwt.AlgorithmValidator("RS256")) != nil {
			c.JSON(http.StatusUnauthorized, "could not validate signing algorithm")
			return
		}

		if jot.Validate(jwt.IssuerValidator("https://pungy.eu.auth0.com/")) != nil {
			c.JSON(http.StatusUnauthorized, "could not validate issuer")
			return
		}

		if jot.Validate(jwt.ExpirationTimeValidator(time.Now())) != nil {
			c.JSON(http.StatusUnauthorized, "token has expired")
			return
		}
	}
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
	authorized.DELETE("/todo", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
