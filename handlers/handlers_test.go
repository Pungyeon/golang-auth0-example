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
	r := gin.New()
	r.GET("/todo", todoHandler.GetTodoListHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todo", nil)
	req.Header.Add("username", "lja")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("GET /todo -> Response[%d]", w.Code)
	}
}
