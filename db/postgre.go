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
	sqlConnection, err := sql.Open("Postgre", connectionString)

	if err != nil {
		panic(err)
	}
	return &PostgreTodoDB{
		connection: sqlConnection,
	}
}

// PostgreConfig holds the configuration for a Postgre connection
type PostgreConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
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

// Get returns a todo from a Postgre database
/* func (pq *PostgreTodoDB) Get(id int) todo.Todo {
	var t todo.Todo
	row := pq.connection.QueryRow(`SELECT * FROM todo WHERE uid=$1`, id)
	err := row.Scan(&t.UUID, &t.Title, &t.Description, &t.Username, &t.Completed)
	if err != nil {
		log.Println(err)
		return t
	}
	return t

} */

// Get will return all todos which are tied to a specified user
func (pq *PostgreTodoDB) Get(user string) []todo.Todo {
	var todos []todo.Todo
	rows, err := pq.connection.Query("SELECT * FROM todo WHERE username=$1", user)
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
	stmt, err := pq.connection.Prepare(`INSERT INTO todo(title, description, username, completed) VALUES($1,$2,$3,$4) returning uid;`)
	if err != nil {
		return "", err
	}
	var uid int
	stmt.QueryRow(
		t.Title, t.Description, t.Username, t.Complete,
	).Scan(&uid)
	return string(uid), err
}

// Put edits a todo in a Postgre database
func (pq *PostgreTodoDB) Put(t todo.Todo) error {
	stmt, err := pq.connection.Prepare("UPDATE todo SET title=$1 description=$2 completed=$3 where uid=$4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Title, t.Description, t.Complete, t.UUID)
	return err
}

// Complete will complete a todo specified by id
func (pq *PostgreTodoDB) Complete(id string) error {
	stmt, err := pq.connection.Prepare("UPDATE todo SET completed=$1 where uid=$2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(true, id)
	return err
}

// Delete removes a todo from a Postgre database
func (pq *PostgreTodoDB) Delete(id string) error {
	_, err := pq.connection.Exec("DELETE FROM todo WHERE uid=$1", id)
	if err != nil {
		return err
	}
	return nil
}
