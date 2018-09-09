# PostgreSQL 

## Preparation
We are going to be using a PostgreSQL as a database. The first thing that we will do, is to create a database and table in our PostgreSQL database. The names of both dtabase and table are not of importance, however, for the rest of this tutorial, all references to both database and table will assume the following names:

* Database: todo
* Table: todos

Personally, I am using docker to run this database on my development machine. I have added a `Dockerfile` which will build the PostgreSQL database and ensure that both databse and table are present. You can also run the official PostgreSQL docker image (or a non-docker deployment) and use pgAdmin (which has a GUI) to ensure that both database and table are created. To download PostgreSQL and pgAdmin:

* PostgreSQL: https://www.postgresql.org/download/
* PgAdmin: https://www.pgadmin.org/download/


## Dependency Injection
// Explain dependency injection

## Preparing our Code for Dependency Injection
First, we will create a new folder `db` with a new file `memory.go`. Which will contain our previous implementation, which was a simple in-memory implementation of a database (essentially). However, we will also create a file called `db.go`, which will implement our database interface.

```go
package db

import "github.com/Pungyeon/golang-auth0-example/todo"

// TodoDB interface for interacting with backend database
type TodoDB interface {
	Get(id int) todo.Todo
	GetAllUserTodos(user string) []todo.Todo
	Add(t todo.Todo) (int, error)
	// Put(t todo.Todo) error // Maybe use this instead of complete?
	Complete(string id)
	Delete(id int) error
}
```

// explain interfaces? might not be necessary depending on the dependency injection explanation
NOTE: You don't need to write the input parameter name. So, essentially `Get(id int)` can be shortened to `Get(int)`. However, I prefer keeping these variable names, to indicate a description of what we are passing to the interface function.

Notice that this interface is implementing all the public methods of our previous implementation... Convenient, nice.

Now, on to our `memory.go` file, which will implement our in-memory database as a class, rather than as a static class. This means that our file `.\todo\todo.go` will become:

```go
package todo

// Todo data structure for task with a description of what to do
type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}
```

Obviously, we are removing (or at least moving) all logic. However, more notably, we are moving our initialisation logic into the `memory.go` file.

```go
package db

import (
	"errors"
	"sync"

	"github.com/Pungyeon/golang-auth0-example/todo"
	"github.com/rs/xid"
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
func (db *InMemoryDB) Add(message string) string {
	t := db.newTodo(message)
	db.mtx.Lock()
	db.list = append(db.list, t)
	db.mtx.Unlock()
	return t.ID
}

// Delete will remove a Todo from the Todo list
func (db *InMemoryDB) Delete(id string) error {
	location, err := db.findTodoLocation(id)
	if err != nil {
		return err
	}
	db.removeElementByLocation(location)
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed
func (db *InMemoryDB) Complete(id string) error {
	location, err := db.findTodoLocation(id)
	if err != nil {
		return err
	}
	db.setTodoCompleteByLocation(location)
	return nil
}

func (db *InMemoryDB) newTodo(msg string) todo.Todo {
	return todo.Todo{
		ID:       xid.New().String(),
		Title:    msg,
		Complete: false,
	}
}

func (db *InMemoryDB) findTodoLocation(id string) (int, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	for i, t := range db.list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func isMatchingID(a string, b string) bool {
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
```

This new implementation of our in-memory database is not that different from our previous implementation. However, as mentioned before, we have moved the intialisation logic to this file, which is now represented by the `NewInMemoryDB` function. This means, that this is no longer automagically initialised, whenever we load the package. We now have to initialise this manually, but don't worry! We will get back to that in just a bit... I promise :* 

We are now going to implement this exact same interface, but rather than it being in-memory, we are going to implement this for our PostgreSQL database (i promised we would be getting back to that). So, let's create a new file in our `db` folder named `postgre.go`. In this file, we will implement the methods of our `TodoDB` interface:

```go
package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Pungyeon/pq-go/todo"
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

// Get returns a todo from a Postgre database
/* func (pq *PostgreTodoDB) Get(id int) todo.Todo {
	var t todo.Todo
	row := pq.connection.QueryRow(`SELECT * FROM todo WHERE uid=$1`, id)
	err := row.Scan(&t.UID, &t.Title, &t.Description, &t.Username, &t.Completed)
	if err != nil {
		log.Println(err)
		return t
	}
	return t

} */

// GetAllUserTodos will return all todos which are tied to a specified user
func (pq *PostgreTodoDB) Get(user string) []todo.Todo {
	var todos []todo.Todo
	rows, err := pq.connection.Query("SELECT * FROM todo WHERE username=$1", user)
	if err != nil {
		log.Println(err)
		return []todo.Todo{}
	}
	for rows.Next() {
		var t todo.Todo
		err = rows.Scan(&t.UUID, &t.Title, &t.Description, &t.Username, &t.Completed)
		if err != nil {
			continue
		}
		todos = append(todos, t)
	}
	return todos
}

// Insert a todo into a Postgre database
func (pq *PostgreTodoDB) Insert(t todo.Todo) (int, error) {
	stmt, err := pq.connection.Prepare(`INSERT INTO todo(title, description, username, completed) VALUES($1,$2,$3,$4) returning uid;`)
	if err != nil {
		return 0, err
	}
	var uid int
	stmt.QueryRow(
		t.Title, t.Description, t.Username, t.Completed,
	).Scan(&uid)
	return uid, err
}

// Put edits a todo in a Postgre database
func (pq *PostgreTodoDB) Put(t todo.Todo) error {
	stmt, err := pq.connection.Prepare("UPDATE todo SET title=$1 description=$2 completed=$3 where uid=$4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Title, t.Description, t.Completed, t.UUID)
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
func (pq *PostgreTodoDB) Delete(id int) error {
	_, err := pq.connection.Exec("DELETE FROM todo WHERE uid=$1", id)
	if err != nil {
		return err
	}
	return nil
}
```

// explain the code

Ok. So now we have implemented the interfaces for `TodoDB` on two different structs. But, now our `handlers.go` file, located in our `handlers` folder is totally broken. We need to fix this. How? Well, one way to do this, is by ensuring that we can initialise this as a struct as well and enabling dependency injection by passing a `TodoDB` on initialisation. Therefore, we need to change our code to the following:

// We are also having to implement the check of a user

```go 
package handlers

import (
	"encoding/json"
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
	c.JSON(http.StatusOK, handler.todo.Get())
}

// AddTodoHandler adds a new todo to the todo list
func (handler *TodoHandler) AddTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	t, err := handler.todo.Add(todoItem.Title)
	if err != nil {
		c.JSON(statusCode, err)
	}
	c.JSON(statusCode, gin.H{"id": t})
}

// DeleteTodoHandler will delete a specified todo based on user http input
func (handler *TodoHandler) DeleteTodoHandler(c *gin.Context) {
	todoID := c.Param("id")
	if err := handler.todo.Delete(todoID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// CompleteTodoHandler will complete a specified todo based on user http input
func (handler *TodoHandler) CompleteTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	if handler.todo.Complete(todoItem.UUID) != nil {
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
```

// explain code 

So that's great, but now our `main.go` file is broken! Time to fix that as well:

```go
package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/auth0-community/auth0"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/Pungyeon/golang-auth0-example/handlers"
	"github.com/Pungyeon/golang-auth0-example/db"
)

/* USAGE:
 * Set the environment variables:
 * 	AUTH0_CLIEN_ID
 * 	AUTH0_DOMAIN: "https://pungy.eu.auth0.com/" (very importantly not omitting the last /)
 */

var (
	audience string
	domain   string
	todoHandler handlers.TodoHandler
)

func main() {
	backend := flag.String("db", "memory", "specify which backend to use for our todo ap: ('memory' | 'postgres')"
	config := flag.String("config", "config.json", "specify the relative filepath of the config to use for the postgres db connecdtion")
	flag.Parse()

	switch *backend {
	case "memory":
		todoHandler = db.NewInMemoryDB()
	case "postgres":
		config, err := db.ReadPostgreconfig(*config)
		if err != nil {
			panci(err)
		}
		todoHandler = db.NewPostgreTodoDB()
	}

	setAuth0Variables()
	r := gin.Default()

	// This will ensure that the angular files are served correctly
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	authorized := r.Group("/")
	authorized.Use(authRequired())
	authorized.GET("/todo", todoHandler.GetTodoListHandler)
	authorized.POST("/todo", todoHandler.AddTodoHandler)
	authorized.DELETE("/todo/:id", todoHandler.DeleteTodoHandler)
	authorized.PUT("/todo", handtodoHandlerlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_CLIENT_ID")
	domain = os.Getenv("AUTH0_DOMAIN")
}

// ValidateRequest will verify that a token received from an http request
// is valid and signy by authority
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
```



