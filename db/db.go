package db

import "github.com/Pungyeon/golang-auth0-example/todo"

// TodoDB interface for interacting with backend database
type TodoDB interface {
	// Get(id int) todo.Todo
	// GetAllUserTodos(user string) []todo.Todo
	Get(user string) []todo.Todo
	Add(t todo.Todo) (string, error)
	// Put(t todo.Todo) error // Maybe use this instead of complete?
	Complete(id string, username string) error
	Delete(id string, username string) error
}
