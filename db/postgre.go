package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Pungyeon/golang-auth0-example/todo"
	_ "github.com/lib/pq"
)

// PostgreTodoDB is a connection to a PostgreSQL todo database
// and holds methods for interacting with Todos in a PostgreSQL
type PostgreTodoDB struct {
	connection *sql.DB
} // @Implements TodoDB

// NewPostgreTodoDB returns a new PostgreTodoDB pointer
func NewPostgreTodoDB(config PostgreConfig) *PostgreTodoDB {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)
	sqlConnection, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}
	return &PostgreTodoDB{
		connection: sqlConnection,
	}
}

// PostgreConfig holds the configuration for a Postgre connection
type PostgreConfig struct {
	DBHost     string `json:"host"`
	DBPort     int    `json:"port"`
	DBUser     string `json:"user"`
	DBPassword string `json:"password"`
	DBName     string `json:"db_name"`
}

// ReadPostgreConfig will attempt to read a config file based on the
// filepath specified, and return a new PostgreConfig
func ReadPostgreConfig(filepath string) (PostgreConfig, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return PostgreConfig{}, err
	}
	var config PostgreConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return PostgreConfig{}, err
	}
	return config, nil
}

// GetByID returns a todo from a Postgre database retrieving the todo via. its UUID
func (pq *PostgreTodoDB) GetByID(id string) todo.Todo {
	var t todo.Todo
	row := pq.connection.QueryRow(`SELECT * FROM todos WHERE uuid=$1`, id)
	err := row.Scan(&t.UUID, &t.Title, &t.Description, &t.Username, &t.Complete)
	if err != nil {
		log.Println(err)
		return t
	}
	return t
}

// Get will return all todos which are tied to a specified user
func (pq *PostgreTodoDB) Get(user string) []todo.Todo {
	var todos []todo.Todo
	rows, err := pq.connection.Query("SELECT * FROM todos WHERE username=$1", user)
	if err != nil {
		log.Println(err)
		return []todo.Todo{}
	}
	for rows.Next() {
		var t todo.Todo
		err = rows.Scan(&t.UUID, &t.Title, &t.Description, &t.Username, &t.Complete)
		if err != nil {
			continue
		}
		todos = append(todos, t)
	}
	return todos
}

// Add a todo into a Postgre database
func (pq *PostgreTodoDB) Add(t todo.Todo) (string, error) {
	stmt, err := pq.connection.Prepare(`INSERT INTO todos(uuid, title, description, username, completed) VALUES($1,$2,$3,$4,$5) returning uuid;`)
	if err != nil {
		return "", err
	}
	var uuid string
	err = stmt.QueryRow(
		t.UUID, t.Title, t.Description, t.Username, t.Complete,
	).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return string(uuid), err
}

// Complete will complete a todo specified by id
func (pq *PostgreTodoDB) Complete(id string, username string) error {
	stmt, err := pq.connection.Prepare("UPDATE todos SET completed=$1 where uuid=$2 AND username=$3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(true, id, username)
	return err
}

// Delete removes a todo from a Postgre database
func (pq *PostgreTodoDB) Delete(id string, username string) error {
	_, err := pq.connection.Exec("DELETE FROM todos WHERE uuid=$1 AND username=$2", id, username)
	if err != nil {
		return err
	}
	return nil
}
