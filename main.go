package main

import (
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Pungyeon/golang-auth0-example/handlers"
)

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

	r.GET("/chat", handlers.GetChat)
	r.POST("/chat", handlers.SendChat)
	r.GET("/all", handlers.AllRooms)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
