package todo

import (
	"github.com/rs/xid"
)

// Todo data structure for task with a description of what to do
type Todo struct {
	UUID        string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Username    string `json:"user"`
	Complete    bool   `json:"complete"`
}

// New returns a new todo
func New(msg string) Todo {
	return Todo{
		UUID:     xid.New().String(),
		Title:    msg,
		Complete: false,
	}
}
