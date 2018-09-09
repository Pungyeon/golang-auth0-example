package db

import (
	"errors"
	"sync"

	"github.com/Pungyeon/golang-auth0-example/todo"
)

// InMemoryDB is an implementation of the TodoDB
// working exclusively in memory. This implementation
// is ephemeral, meaning data will onyl be stored
// for the duration of the program and deleted on program exit
type InMemoryDB struct {
	list []todo.Todo
	mtx  sync.RWMutex
} // @Implements TodoDB

// NewInMemoryDB will return a newly intialised in-memory database
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		list: []todo.Todo{},
	}
}

// Get retrieves all elements from the todo list
func (db *InMemoryDB) Get(user string) []todo.Todo {
	todoListForUser := []todo.Todo{}
	for _, todo := range db.list {
		if todo.Username == user {
			todoListForUser = append(todoListForUser, todo)
		}
	}
	return todoListForUser
}

// Add will add a new todo based on a message
func (db *InMemoryDB) Add(t todo.Todo) (string, error) {
	db.mtx.Lock()
	db.list = append(db.list, t)
	db.mtx.Unlock()
	return t.UUID, nil
}

// Delete will remove a Todo from the Todo list
func (db *InMemoryDB) Delete(id string, username string) error {
	location, err := db.findTodoLocation(id, username)
	if err != nil {
		return err
	}
	db.removeElementByLocation(location)
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed
func (db *InMemoryDB) Complete(id string, username string) error {
	location, err := db.findTodoLocation(id, username)
	if err != nil {
		return err
	}
	db.setTodoCompleteByLocation(location)
	return nil
}

func (db *InMemoryDB) findTodoLocation(id string, username string) (int, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	for i, t := range db.list {
		if isMatching(t.UUID, id) && isMatching(t.Username, username) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id and username")
}

func isMatching(a string, b string) bool {
	return a == b
}

func (db *InMemoryDB) removeElementByLocation(i int) {
	db.mtx.Lock()
	db.list = append(db.list[:i], db.list[i+1:]...)
	db.mtx.Unlock()
}

func (db *InMemoryDB) setTodoCompleteByLocation(location int) {
	db.mtx.Lock()
	db.list[location].Complete = true
	db.mtx.Unlock()
}
