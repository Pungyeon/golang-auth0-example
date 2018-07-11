package todo

import "testing"

func TestAddTodo(t *testing.T) {
	initial := len(list)
	Add("i am a new todo")
	if initial+1 != len(list) {
		t.Error("add does not add another todo message to the list")
	}
}

func TestDeleteTodo(t *testing.T) {
	id := Add("please delete me")
	err := Delete(id)
	if err != nil {
		t.Error(err)
	}

	if l, err := findTodoLocation(id); err == nil {
		t.Errorf("todo not deleted, still able to find via. id: %d", l)
	}
}

func TestCompleteTodo(t *testing.T) {
	id := Add("please complete me")
	Complete(id)

	location, err := findTodoLocation(id)
	if err != nil {
		t.Error(err)
	}

	if list[location].Complete == false {
		t.Errorf("did not complete todo: %v", list[location])
	}
}

func TestGetTodos(t *testing.T) {
	Add("please complete me")
	Add("please complete me")
	Add("please complete me")
	Add("please complete me")

	retrievedList := Get()
	if len(retrievedList) != len(list) {
		t.Error("getting todo list returns incorrect result")
	}
}

func TestDeleteFailure(t *testing.T) {
	if Delete("lakjsdfl") == nil {
		t.Error("deletion of non-existent item, did not return 404 error")
	}
}

func TestCompleteFailure(t *testing.T) {
	if Complete("lakjsdfl") == nil {
		t.Error("completion of non-existent item, did not return 404 error")
	}
}
