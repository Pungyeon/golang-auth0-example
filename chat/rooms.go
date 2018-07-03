package chat

import (
	"fmt"
	"sync"
)

// Room is a struct for defining a chat room that is defined
// by name and the messages bound to the room
type Room struct {
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
}

// Message is a struct for defining chat messages
// that are bound to a user (name) and its content
type Message struct {
	Room    string `json:"room"`
	User    string `json:"user"`
	Content string `json:"content"`
}

var (
	rooms map[string]*Room
	mtx   sync.RWMutex
	once  sync.Once
)

func init() {
	rooms = map[string]*Room{}
}

// GetRooms retrieves the rooms static object
func GetRooms() *map[string]*Room {
	return &rooms
}

// AddMessage will add a message to the relevant chat room
func (r *Room) AddMessage(msg Message) {
	room := GetRoom(msg.Room)
	mtx.Lock()
	room.Messages = append(room.Messages, msg)
	mtx.Unlock()
}

// GetRoom retrieves specific room, specified by room name
func GetRoom(name string) *Room {
	mtx.RLock()
	r, ok := rooms[name]
	mtx.RUnlock()
	if !ok {
		fmt.Println("creating new room")
		nr := Room{Name: name, Messages: []Message{}}
		mtx.Lock()
		rooms[name] = &nr
		mtx.Unlock()
		return rooms[name]
	}
	return r
}
