package db

import (
	"testing"

	"github.com/Pungyeon/golang-auth0-example/todo"
)

func TestAddTodo(t *testing.T) {
	todos := NewInMemoryDB()
	todo := todo.New("i am a new todo", "lja")
	todos.Add(todo)
	if len(todos.list) != 1 {
		t.Error("add does not add another todo message to the list")
	}
}

func TestDeleteTodo(t *testing.T) {
	todos := NewInMemoryDB()
	todo := todo.New("i am a new todo", "lja")
	id, err := todos.Add(todo)
	if err != nil {
		t.Error(err)
	}
	err = todos.Delete(id)
	if err != nil {
		t.Error(err)
	}

	if l, err := todos.findTodoLocation(id); err == nil {
		t.Errorf("todo not deleted, still able to find via. id: %d", l)
	}
}

func TestCompleteTodo(t *testing.T) {
	todos := NewInMemoryDB()
	todo := todo.New("i am a new todo", "lja")
	id, err := todos.Add(todo)
	if err != nil {
		t.Error(err)
	}
	todos.Complete(id)

	location, err := todos.findTodoLocation(id)
	if err != nil {
		t.Error(err)
	}

	if todos.list[location].Complete == false {
		t.Errorf("did not complete todo: %v", todos.list[location])
	}
}

func TestGetTodos(t *testing.T) {
	todos := NewInMemoryDB()
	todo := todo.New("i am a new todo", "lja")
	todos.Add(todo)
	todos.Add(todo)
	todos.Add(todo)
	todos.Add(todo)

	retrievedList := todos.Get("lja")
	if len(retrievedList) != len(todos.list) {
		t.Error("getting todo list returns incorrect result")
	}
}

func TestDeleteFailure(t *testing.T) {
	todos := NewInMemoryDB()
	if todos.Delete("lakjsdfl") == nil {
		t.Error("deletion of non-existent item, did not return 404 error")
	}
}

func TestCompleteFailure(t *testing.T) {
	todos := NewInMemoryDB()
	if todos.Complete("lakjsdfl") == nil {
		t.Error("completion of non-existent item, did not return 404 error")
	}
}
