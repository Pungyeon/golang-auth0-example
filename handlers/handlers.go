package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Pungyeon/golang-auth0-example/db"
	"github.com/Pungyeon/golang-auth0-example/todo"
	"github.com/gin-gonic/gin"
)

// TodoHandler contains all the API endpoints for our todo application
type TodoHandler struct {
	todo db.TodoDB
}

// NewTodoHandler will return a new TodoHandler, specifying the type
// of backend (database) to use, via. the db input parameter
func NewTodoHandler(db db.TodoDB) *TodoHandler {
	return &TodoHandler{
		todo: db,
	}
}

// GetTodoListHandler returns all current todo items
func (handler *TodoHandler) GetTodoListHandler(c *gin.Context) {
	// get user
	username, err := getUsername(c.Request)
	if err != nil {
		terminateWithError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, handler.todo.Get(username))
}

func getUsername(r *http.Request) (string, error) {
	username := r.Header.Get("username")
	if username == "" {
		return username, errors.New("could not retrieve username from header, suggesting it hasn't been added do to a previous error")
	}
	return username, nil
}

// AddTodoHandler adds a new todo to the todo list
func (handler *TodoHandler) AddTodoHandler(c *gin.Context) {
	username, err := getUsername(c.Request)
	if err != nil {
		terminateWithError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	todoItem.Username = username
	t, err := handler.todo.Add(todoItem)
	if err != nil {
		c.JSON(statusCode, err)
	}
	c.JSON(statusCode, gin.H{"id": t})
}

// DeleteTodoHandler will delete a specified todo based on user http input
func (handler *TodoHandler) DeleteTodoHandler(c *gin.Context) {
	username, err := getUsername(c.Request)
	if err != nil {
		terminateWithError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	todoID := c.Param("id")
	if err := handler.todo.Delete(todoID, username); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// CompleteTodoHandler will complete a specified todo based on user http input
func (handler *TodoHandler) CompleteTodoHandler(c *gin.Context) {
	username, err := getUsername(c.Request)
	if err != nil {
		terminateWithError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	if handler.todo.Complete(todoItem.UUID, username) != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (todo.Todo, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return todo.Todo{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (todo.Todo, int, error) {
	var todoItem todo.Todo
	err := json.Unmarshal(jsonBody, &todoItem)
	if err != nil {
		return todo.Todo{}, http.StatusBadRequest, err
	}
	return todoItem, http.StatusOK, nil
}
