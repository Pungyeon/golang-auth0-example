package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Pungyeon/golang-auth0-example/chat"
	"github.com/gin-gonic/gin"
)

func AllRooms(c *gin.Context) {
	c.JSON(http.StatusOK, chat.GetRooms())
}

// GetChat will respond with all recorded chat messages
func GetChat(c *gin.Context) {
	roomName := c.Query("room")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room name not specified, please specify a room name"})
		return
	}

	c.JSON(http.StatusOK, chat.GetRoom(roomName))
}

// SendChat will add a new chat messages to the chat
func SendChat(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer c.Request.Body.Close()

	var message chat.Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	chat.GetRoom(message.Room).AddMessage(message)
	c.JSON(http.StatusOK, chat.GetRooms())
}
