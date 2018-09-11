package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pungyeon/golang-auth0-example/db"
	"github.com/gin-gonic/gin"
)

func TestAddAndGetHandler(t *testing.T) {
	// SETUP
	todoHandler := NewTodoHandler(
		db.NewInMemoryDB(),
	)
	/*
		td := todo.New("this is a test", "lja")
		data, _ := json.Marshal(td)
	*/
	r := gin.Default()
	r.GET("/todo", todoHandler.GetTodoListHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todo", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("response from GET -> /todo was not 200, something went wrong")
	}
}
