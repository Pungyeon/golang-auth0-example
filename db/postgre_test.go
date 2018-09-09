package db

import (
	"testing"

	"github.com/Pungyeon/golang-auth0-example/todo"
)

var (
	config = PostgreConfig{
		DBHost:     "172.16.1.68",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "todo",
	}
)

func TestPostgresAddAndGet(t *testing.T) {
	pq := NewPostgreTodoDB(config)
	newtodo := todo.New("hello there", "lja")
	id, err := pq.Add(newtodo)
	if err != nil {
		t.Error(err)
	}

	td := pq.GetByID(id)
	if td.Title != "hello there" {
		t.Error("could not fetch the newly added todo by id")
	}

	tdlist := pq.Get("lja")
	if len(tdlist) == 0 {
		t.Error("Get() returns an empty list, despite list being demonstrably non-empty")
	}
}

func TestPostgresDelete(t *testing.T) {
	pq := NewPostgreTodoDB(config)
	newtodo := todo.New("hello there", "lja")
	id, err := pq.Add(newtodo)
	if err != nil {
		t.Error(err)
	}

	err = pq.Delete(id, "lja")
	if err != nil {
		t.Error(err)
	}

	td := pq.GetByID(id)
	if td.UUID != "" {
		t.Error("todo still exists after attempting to delete them")
	}
}

func TestPostgresComplete(t *testing.T) {
	pq := NewPostgreTodoDB(config)
	newtodo := todo.New("hello there", "lja")
	id, err := pq.Add(newtodo)
	if err != nil {
		t.Error(err)
	}
	err = pq.Complete(id, "lja")
	if err != nil {
		t.Error(err)
	}
	td := pq.GetByID(id)
	if err != nil {
		t.Error(err)
	}
	if td.Complete != true {
		t.Error("completed todo item is not set to complete in the database")
	}
}

func TestReadPostgreConfig(t *testing.T) {
	config, err := ReadPostgreConfig("../config.test.json")
	if err != nil {
		t.Error(err)
	}
	if config.DBUser != "postgres" {
		t.Error("config could not be loaded, returning wrong config values")
		t.Error("expected: postgres, actual: " + config.DBUser)
	}
}
