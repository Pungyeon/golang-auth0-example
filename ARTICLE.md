# Building Applications with Golang and Angular

**TL;DR:** In this article, you will learn how to build modern applications with Golang and Angular. Throughout the article, you will build a secure Golang API that will support a ToDo list application that you will develop with Angular. To facilitate the identity management, you will integrate Auth0 both in your backend and in your frontend.

## Why Choosing Golang and Angular
The technologies of choice for this article, is golang for the backend and angular for the frontend. But why have we chosen these two technologies? 

### Golang
Golang (or Go), is Google's very own programming language. It's a statically typed and compiled language. Over the last few years Golang has become a very popular language, being the language of choice for projects such as Docker, Kubernetes and all of Hashicorps suite of programs. Golang is not quite like other programming languages. It has a very strong standard library and can get very far before having to use 3rd party libraries, it compiles to a single binary, has amazing support for concurrency, testing & benchmarking and has a fantastic community behind it.

Golang is by no means the perfect language, but the reason I use it every day, is for the simple reason that it get's the job done. Other than Python or Ruby, I don't believe there is a language as effective in time-to-deliver as Golang. The reason why I haven't chosen Ruby or Python being, I prefer statically typed compiled languages. It results in a faster application and avoids dependency hell for the end user. Golang also has excellent cross-platform compilation support, making it even more attractive as a programming language.

### Angular
Currently, the only real viable choice for writting a frontend is using Javascript. All browsers support it to a satisfactory degree and so we avoid strange dependency and compatibility issues. So, that makes choosing the frontend technology easy... right? Nope! As of writing this article, there are currently three major frameworks for writing web applications: React, Vue & Angular (and 100 other smaller frameworks). Now, you can 100% find an article which states that "Framework X is much better than frameworks Y & Z, because blah" or an article describing that all three frameworks mentioned are essentially dead, because of some new framework. I'm not going to go into the politics of frontend frameworks, it's just not worth the effort. My own opinion is, that all three frameworks accomplish the exact same thing and the end result that they produce are... quite similar. All use yarn or npm to download 100+ MB's of god knows what in node_modules, and then compress all of this into  smaller html and js files. 

I like Angular, because it comes with everything that I expect to come with a web framework. Some standard components, easy management of libraries, routing and authentication support. It makes life easy, and I like life, when it's easy. React is also a super strong framework, but I have OCD issues with javascript in my html, just like many people have issues with their different foods touching one another on their plate. Lastly, we have Vue, which is also a very strong framework. As of now, I cannot tell you a substantial reason to choose Angular over Vue, nor the other way around. However, I started out with Angular and Angular does the job for me.

Either way, Angular is another Google backed project, much like Golang. It has sprouted from it's old brother AngularJS, which was **THE** javascript framework a few years ago. The new Angular is not as popular, which I think there are a few reasons for. There is more competition with React and Vue. There were a lot of breaking changing at the start of Angular and the last reason, which cause a lot of confusion is the versioning convention. Instead of making versions such as 2.2.1.1, it was decided that all breaking changes would cause a total version aggregation. So, the current version of Angular (as of writing this article) is Angular 6. This makes it difficult to Google (ironically) and also weird to refer to, and there are many different ways to refer to it: Angular, Angular 2+, Angular 6, etc.

Despite all of this, I still think that Angular is an excellent frontend framework for JavaScript. I love that it forces uses to use TypeScript and thereby standardising the structure of code and as mentioned earlier, I love the angular-cli toolbox, which comes with all the tools I need and expect out-of-the-box.

## Prequisites for Golang and Angular
### Golang
We need to install Golang. That's easy. Golang is awesome in that way. For installation instructaions, please visit https://golang.org/doc/install

### Angular
First we need to install npm and node. That can be done using these instructions https://nodejs.org/en/download/. From here on, we can type this command in our terminal or command line: 

> npm install -g @angular/cli

Boom, we are ready.

## Building the Golang API
Now we are going to build our Golang API. We will be using the web server framework 'Gin' for this. Gin is, like many other go frameworks, an open source project which simplifies creating API endpoints. Keep in mind, that nothing we will be building in this article is impossible to do with the standard library of go. The only reason we are using gin, is because it simplifies and standardises our process a little, making life easier. We like life, when life is easy.

### Creating our In-Memory ToDo List
Before we start writing our web server, we will start writing our component for handling a Todo list. Our implementation will be a static object, which will store all todo items in-memory and perform CRUD operations on these todo items. Essentially, we are mocking a very simple database. Typically, this is not a bad way to start out development; Implementing a mock version of your database, before implementing your actual database. Not only does it makes testing easier (and something that you can do from the beginning of your project), but it also helps implying an interface for our store (or database). 

So let's get started with our project. Golang by default will look for pacakges in the GO_PATH environment variable, which is places in the user directory of the system (i.e on unix systems ~/go & on windows %USERPROFILE%/go). Packages are then stored in ~/go/src/, and therefore, placing our projects here, will make our lives a lot easier (remember, we like this). I have placed mine in the directory ~/go/src/github.com/Pungyeon/golang-auth0-example and will refer to this directory as root (or './') from here on.

In our root directory, create a new folder named 'store' and in this folder place a new file called store.go. In this file, we will write the following content:

#### ./store/store.go
```go
package todo

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	list []Todo
	mtx  sync.RWMutex
	once sync.Once
)

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []Todo{}
}

// Todo data structure for task with a description of what to do
type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

// Get retrieves all elements from the todo list
func Get() []Todo {
	return list
}

// Add will add a new todo based on a message
func Add(message string) string {
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}

// Delete will remove a Todo from the Todo list
func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed
func Complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompleteByLocation(location)
	return nil
}

func newTodo(msg string) Todo {
	return Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func isMatchingID(a string, b string) bool {
	return a == b
}
```

Explaining everything from top to bottom, we begin with our global variables for this package:
* list -> Is a type of a Todo array 
* mtx is our global mutex for all variables of this package. 
* once contains a golang native functionality `sync.Once`, which will ensure an operation only is run once. 

Our `init()` will run our `initialiseList()`, but ensure that it is only run once. The `init()` function is another golang native function, which is run on package initialisation (whenever the package is loaded). The `initialiseList()` function will reset / initialise our todo list, so we want to make sure that this is only executed once per runtime.

Next, we write our Todo struct, which is the base of our store. If other packages used this struct, we would place it, in another package for itself. But for this simple application, placing it here, will suffice. The struct defines the id, the contained message and whether the todo item is complete or active. We also map all properties of our struct to it's json equivalent. This is a very useful feature in go. The struct property naming convention in go, makes all properties starting with a capital letter public and all starting with a small letter private. In json, usually all properties of an object starts with a non-capital letter, and therefore this mapping ensures we can stick to both naming conventions.

```
This feature is also exceptionally useful when mapping json to a struct. Let's say that you don't like the naming convention used in the json, or it represents something other than your data, you can easily map the json key to a different property name, using this method.
```

Next, we will implement the first method of our todo store: `Get()`. This method is capitalised and is therefore public, meaning it can be accessed by other packages. Get very simply returns our current static todo list (the global list variable).

The next method is `Add()`, which will create a new Todo based on a user input message and then append this todo to our list. Notice that we are using our mutex to `Lock()` before we append to our list and then `Unlock()` again once this operation has ended. This is very important, as we might send multiple operations at the same time. If these operations try to access the same memory, we can run into a race-condition. This would be bad, and would actually cause golang to crash. To avoid this, we use the mutex, which is scoped to our package.

The third public method is `Delete()` which will remove an item from the list. The last public method one is `Complete()` which will mark a todo item as complete in our list, based on it's `id`. 

```
Notice that all our public functions are very simply. They aren't doing much more than a single operation. Finding an item and then either deleting or completing them. This makes our code really easy to read and easy to understand. These principles come from 'Clean Code' written by Robert C. Martin, whose material I can whole-heartedly recommend!
```

Now, we dive deep into the private functions of our package. Our `newTodo` function will take in a message in form of a string, and return a new Todo object, with a newly created UUID (in form of a string) and the complete flag set to false. The next function `findTodoLocation`, will take an input `id` and find a todo item, with a matching `id`. If there are no matches found after iterating over all the items in our list, we will return an error saying that we didn't find any items. Notice, that we are using our mutex again. This time, we are just using the `RLock` function, since we will only be reading our list and not writing to it and `RLock` is slightly more efficient than `Lock`.

Be aware, that we are using a third party package, `xid`, for generating our GUID. To obtain this package before compiling our application, we will have to run the following command:

> go get github.com/rs/xid

Next, we have our `removeElementByLocation` function. We are setting our list variable to a new list, which contains all elements of our list up to a given location, appended with all elements after (but not including) the same given location. This means, we can give our function a location and it will set our list to an new list, without that given location, essentially deleting it from our list.

The very last function in our store packate is our `setTodoCompleteByLocation`. Just like our remove function, this function takes in a location in form of an integer. However, this function is much less complex and simply sets the given location in our list to true.

Putting this together, we can see how our public functions `Delete` and `Complete` work. They find a todo item based on id and returns a location of where that item resides in our list, or returns an error if no item exists. If an error was returned, we pass on this error to the function caller, if not, we perform either a delete or complete operation on that item's location in our list.

### Initilialising our Web Server & Serving Static Files
Now to the web server!

As mentioned earlier, we will write this package in gin. This package can be obtained by writing the following command:

> go get github.com/gin-gonic/gin

Great. Now, in our root folder, we will create our main file (wuhuuu), main.go and write th following code:

#### ./main.go
```go
func main() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	r.GET("/todo", handlers.GetTodoListHandler)
	r.POST("/todo", handlers.AddTodoHandler)
	r.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	r.PUT("/todo", handlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
```

First we instatiate our gin server using `gin.Default()`, this will return an object, in which we can configure and run our web server from. Then we do something, which is a little hacky, so let's explain. Routing in gin is quite specific and I cannot have ambiguous routes for the root path (you can but... you shouldn't). Essentially, gin will complain if you have something matching a route such as `/*`, because this will interfere with every other route in our web server, which will then never be called. In Nodejs (for example), we are able to do this, because the path routing is determined by most specific -> least specific. So a route such as `/api/something` will have precedence over `/*` and this is not the case by default in gin. However, we implement this in our server, by creating a `NoRoute`, matching all routes that have not been specified already. This route function will assume that this call is asking for a file and attempt to find this file. 

If the client asks for a the root path, or if the file is not found, we will serve them 'index.html' (which will be produced from our angular project at a later point). Otherwise, we will serve the client, the file they requested. There are other ways to do this and depending on what you want to achieve, better ways to achieve that. However, for the purpose of this tutorial, this will do just fine. 

Now, we can add our routes to fetch the data from our todo list. They are all pointing to the same path '/todo', but are all using different methods:
* GET: Retrieves our entire todo list
* POST: Adds a new todo to the list
* DELETE: Deletes a todo from the list based on an id 
* PUT: Will change a todo item from uncomplete to complete

Each of our endpoints will be structed in the same manner `r.<METHOD>(<PATH>, <Gin function>)`. Our gin function is basically any function that takes the parameter of a gin.Context pointer. If you go look at the `NoRoute` function, you will see an example of a anonymous function, with the input of a gin.Context pointer.

Quite simple. The last thing our main function does, is to run our web server on port 3000 and panic if an error occurs, while running the web server.

Of course, we are not quite finished, because our handlers don't actually exist. So, we will need to implement them, before we can start our web server.

### Developing the API Endpoints
So, we are going to create a new folder named `handlers`, in this folder, we will write a new file called `handlers.go`. In this file, we will write the code for the implementation of our api endpoints (GET, POST, PUT, DELETE -> for '/todo'). 

Since we have already implemented all this functionality in our todo package, this will be a relatively simple exercise. The final product will look as such:

#### ./handlers/handlers.go
```go
package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Pungyeon/golang-auth0-example/todo"
	"github.com/gin-gonic/gin"
)

// GetTodoListHandler returns all current todo items
func GetTodoListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, todo.Get())
}

// AddTodoHandler adds a new todo to the todo list
func AddTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": todo.Add(todoItem.Message)})
}

// DeleteTodoHandler will delete a specified todo based on user http input
func DeleteTodoHandler(c *gin.Context) {
	todoID := c.Param("id")
	if err := todo.Delete(todoID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// CompleteTodoHandler will complete a specified todo based on user http input
func CompleteTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	if todo.Complete(todoItem.ID) != nil {
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
As mentioned earler, all of our handler functions take in the parameter of a gin.Context pointer. This variable essentially contains the http.Request reader and our http.ResponseWriter writer (which are used as input parameters for http handlers in the standard library). It also contains a lot of metadata about these requests, and helpful functions, making it easier for us to process the data, both incoming and outgoing.

The majority of this code is structed as such:
1. Grab input and convert if necessary
2. Check for errors
3. Perform operation
4. Return error or status ok

At the bottom of the code, I have added to helper functions, specific to parsing input. `convertHTTPBodyToTodo` will read the body from the request, and return it as a `Todo` object. This is done by using the `ioutil.ReadAll` which will read all bytes from an `io.Reader` stream. Once all the bytes have been read, we will convert them from JSON (which is the format the request has sent them in) to a Todo object. This is the code seen in `convertJSONBodyToTodo`. We use the standard library `json.Unmarshal` to try and parse a json object to a specified `interface{}` in this case a Todo. Of course, if we fail at doing so, we return an error.

With these two functions written, it's pretty easy to keep our handler logic really simple and neat. The only other actions we use are `c.JSON` using our gin.Context, to return a response and `gin.H` to create a JSON return. In our `DeleteTodoHandler` function, we are also using c.Param. This function will take the parameter that we specified in our routing in our main function and return a string. Basically, if  our client requests `DELETE /todo/ID123`, this will result in our function extracting `ID123` as an id string and our `todo.Delete` function deleting the todo with this id (assuming that it exists).

### Securing the Golang API with Auth0
So, our web server is working as intended now. We can add, edit, complete and delete our todo list via. our API. I suggest that you try this out, by going to the root and running:

> go run main.go

Which will start our web server on `localhost:3000`. You can then use curl (or a program like Postman) to fiddle around with the different actions of our web server. As an example:

Add a new todo:
> curl POST localhost:3000/todo -d '{"message": "finish writing the article"}'

Get all our current todos:
> curl GET localhost:3000/todo

If this works, we will now get a response with our single todo item on our list, suggesting that I should finish writing this article. Great! Time to celebrate! Not quite... we have one issue. Right now, anyone can use our todo list. That's not good, we want to make sure that only people that we trust can access and edit our todo list. To do this, we will utilise Auth0 as a backend authentication service. First step is to setup a free account on Auth0.

<auth0-instructions></auth0-instructions>

Now that our account has been setup, we need to make some changes on our backend. Essentially, we want our web api to check a service every time that a request is made towards. To do this, we will change our `main.go` file to the following:

#### ./main.go
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
)

/* USAGE:
 * Set the environment variables:
 * 	AUTH0_CLIEN_ID
 * 	AUTH0_DOMAIN: "https://pungy.eu.auth0.com/" (very importantly not omitting the last /)
 */

var (
	audience string
	domain   string
)

func main() {
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
	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddTodoHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

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

In our main function we have added a routing group, called authorized. We then tell gin that this routing group should use a middleware function called `authRequired` and then change all our todo operations, to make them a part of this group, rather than being directly attached to our gin router.

We have also written a middlware function. A middleware function is basically a function that is run, before our actual handler. So, a request towards `GET '/todo'` is no processed in the `authRequired` function and then the `handlers.GetTodoListHandler`.

Another addition is two new global variables audience and domain, which we will need for our authentication against Auth0. These will be retrieved from our environment variables on start, using the function `setAuth0Variables`.

The function `authRequired` is a middleware function. In gin terms, this must return a gin.HandlerFunc, which contains a `Next()` invokation in the body. Basically, our function validates a token which is provided in the incoming request 'Authorization' header. We do this using JWKS (JSON Web Key Set). Essentially, JWKS is a method for verifying JWT, using a public/private key infrastructure. In out case, both the private and public keys are provided by Auth0, so we don't have to do any additional work.

```
To read about JWKS: https://auth0.com/docs/jwks
```

Using the auth0 golang library, makes this extremely simple. All we have to do, is four lines of code to validate our incoming token. If this results in an error we will terminate the current connection, responding to the incoming request with a `http.StatusUnauthorized` (401) and terminating the connection. If the token is validated, then we will send the request onto the next function in the handler chain. 

So let's make sure that this works, by spinning up our server:

> go run main.go

... and then without getting a token, we will send a request to our API expecting a 401 unauthorized:

> curl GET localhost:3000/todo

If this results in a 200, we have done something wrong. If not, awesome. Then it is time to create our UI!

## Developing a ToDo List with Angular
Now that we have our backend sorted, we will proceed with creating a frontend. As stated earlier, this will consist of simple home page with a button to redirect us to the todo list. To access to the todo list, we must be authenticated. Let's get going!

### Initialising the Angular Project
Our new angular project, will be placed in the folder `ui`. This folder will be auto-created on initialisation. So, to initialise the project, use the following command:

> ng new ui

This will place a new angular quickstart project in a new folder: ui. Now, we need to go into our ui folder and download all of our node modules, which are dependencies for our angular project. To do this, run the following command:

> npm install

Last thing we need to do, before we start writing our application, is to add a few CDN links to our index.html file. All of these are also possible to get from `npm install`, or to download locally, but this is the simplest solution as of right now. So, edit the file `./ui/src/app/index.html` to the following:


#### ./ui/src/index.html
```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Auth0 Golang Exaple</title>
  <base href="/">

  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css">
  <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.1.0/css/all.css">
  <script src="https://cdn.auth0.com/js/auth0/9.5.1/auth0.min.js"></script>
</head>
<body>
  <app-root></app-root>
</body>
</html>
```

All we have done is change the title to "Auth0 Golang Example" and added three external CDN dependencies:
1. Bootstrap CSS - to make everything look pretty
2. Font Awesome - to have cool icons for buttons and such
3. Auth0 - Which is a JS library for using Auth0 authentication

Another change that we would like to do for preparation, is to edit the environment files found in `./ui/src/environments`. These files will act as global variables for our application and differentiate between whether we are running a local dev environment, or we are running in production. Change the `environment.prod.ts` to:

#### ./ui/src/environments/environment.prod.ts
```js
export const environment = {
  production: true,
  gateway: '',
  callback: 'http://localhost:3000/callback'
};
```

And if you want to develop using `ng serve`, then change the `environment.ts` (which is considered our dev environment by default), to the following:

#### ./ui/src/environments/environment.ts
```js
export const environment = {
  production: false,
  gateway: 'http://localhost:3000',
  callback: 'http://localhost:4200/callback'
};
```

```
NOTE: The command "ng serve" will start a local web server on port 4200. This web server will detect whenever changes are made to the code base, transpile and reload to serve our angular project with the new changes. This will only transpile the changes, so therefore this is much faster than having to rebuild our project for every change and therefore preferred when developing.
```

### Creating the Welcome & Todo Page
First, we will make a new Angular component. We will quickstart this, using the angular-cli:

> ng g c home

Essentially, this is the shortened version of:

> ng generate component home

This will create a new folder in our `app` folder, named `todo` and place four files in the folder:
* home.component.css - for styling our html
* home.component.html - the html for our page
* home.component.spec.ts - testing file (we won't be using this, you are free to delete it)
* home.component.ts - the TypeScript file for all the javascript to support the html page

```
NOTE: The CLI will also automatically add this class to our `app.module.ts`, in the `@NgModule.declarations` section.
```

Our welcome page, will be extremely simple and we will only have to touch our `home.component.html` file and change it to this:

#### ./ui/src/app/home/home.component.html
```html
<div class="c-block">
  <h3> Welcome Home!</h3>
  <a routerLink="/todo">
    <button class="btn-primary"> Show Todo List </button>
  </a>
</div>
```

All we are doing is creating a title with a link (with a nested button) that will redirect the user to `/todo`. We haven't determined the routing yet, so that won't work right now. We will get to that soon.

However, before we do, let's create a todo component in the same way as our home component:

> ng g c todo

### Developing the ToDo Page & Service
Our todo functionality will be split into two, a component (the one we created in the last section) and a service. The service will be in charge of communicating with the backend via. HTTP. This service may then be used by our component, for ease of use, to communicate with our backend and display the correct information retrieved by the service.

So, let's begin by creating our service. Best practice tells us to create a new folder called `service` and place a new file in there named `todo.service.ts`. This file will consist of the following code: 

#### ./ui/src/app/service/todo.service.ts
```js
import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { environment } from "../../environments/environment";

@Injectable()
export class TodoService {
    constructor(private httpClient: HttpClient) {}

    getTodoList() {
        return this.httpClient.get(environment.gateway + '/todo');
    }

    addTodo(todo: Todo) {
        return this.httpClient.post(environment.gateway + '/todo', todo);
    }

    completeTodo(todo: Todo) {
        return this.httpClient.put(environment.gateway + '/todo', todo);
    }

    deleteTodo(todo: Todo) {
        return this.httpClient.delete(environment.gateway + '/todo/' + todo.id);
    }
}

export class Todo {
    id: string;
    message: string;
    complete: boolean;
}
```

Here we are stating that we want to `export` our TodoService class, thereby making it available for other components to use. On initialisation of our class, we expect an HttpClient as input, which will be available for usage of our class (under the property name `this.httpClient`). This pattern is called dependency injection... but that is a topic for another article. 

```
NOTE: the @Injectable() decorator ensures that we can inject the HttpClient. Without this, you will get some strange console errors.
```

At the bottom of the file, we are exporting another class `Todo` which is a mirror of our todo class from our backend. We use this class throughout our project, to ensure a more accurate description in code, what we are sending and retrieving from our backend. The rest of the service is some very basic HTTP calls to our backend. Notice that we are using the `environment.gateway` variable to determine where our http calls are headed. This makes it easier to change environments.

```
NOTE: This service must also be added to our app.module.ts file, just like all our components. However, since this is a service our components use, we must add it to the 'providers' section in @NgModule. We will also need to add HttpClientModule to the imports section. There is a reference to where and how, at a later point in the tutorial
```

Now, let's use our new todo service, for usage in our todo component. First, let's create our TypeScript logic in the `todo.component.ts` file:

#### ./ui/src/app/todo/todo.component.ts
```js
import { Component, OnInit } from '@angular/core';
import { TodoService, Todo } from '../service/todo.service';

@Component({
  selector: 'app-todo',
  templateUrl: './todo.component.html',
  styleUrls: ['./todo.component.css']
})
export class TodoComponent implements OnInit {

  activeTodos: Todo[];
  completedTodos: Todo[]
  todoMessage: string;

  constructor(private todoService: TodoService) { }

  ngOnInit() {
    this.getAll();
  }

  getAll() {
    this.todoService.getTodoList().subscribe((data: Todo[]) => {
      this.activeTodos = data.filter((a) => !a.complete);
      this.completedTodos = data.filter((a) => a.complete);
    });
  }

  addTodo() {
    var newTodo : Todo = {
      message: this.todoMessage,
      id: '',
      complete: false
    }
    this.todoService.addTodo(newTodo).subscribe(() => {
      this.getAll();
      this.todoMessage = '';
    });
  }

  completeTodo(todo: Todo) {
    this.todoService.completeTodo(todo).subscribe(() => {
      this.getAll();
    });
  }

  deleteTodo(todo: Todo) {
    this.todoService.deleteTodo(todo).subscribe(() => {
      this.getAll();
    })
  }
}
```

Starting from the `export class TodoComponent`, what we have is three properties, which respectively contain our active todo list, our completed todo list and a user input string. Our constructor of our component will expect a TodoService object and store this as a private property todoService.

The first function of our component, `ngOnInit` is a built-in standard of angular and derives from the interface `OnInit`. Essentially and implementation on `OnInit`, will wait for the component to be loaded, before executing the `ngOnInit` function. Essentially, we are waiting to retrieve data, until our component has loaded successfully. In our function, we are executing the method `getAll`. This method, will invoke the `todoService.getTodoList` function. As we know, this function is an http call to our backend to get all of our todo items. The HttpClient in angular does not return with the response, but rather with an Observable. Essentially, this is an object, which we can 'subscribe' to and then whenever new data is retrieved, we will be notified and can hereby respond in some way or another. In this scenario, we are using the observable more like a callback, but I strongly recommend reading up on RxJs and Observables. They are super useful in modern Javascript.

Anyway... the `getAll` function subscribes to data from the `todoService.getTodoList` and whenever data is received, will assigned all active todos to our `activeTodos`, by filtering out any complete items, and do the opposite for our `completedTodos` property. The rest of our class methods are corresponding to our todo service, which in turn was mapped up against our backend api. In other words, we have an add, complete and delete, and whenever an operation is performed, we update our todo list, by retrieving the data again with `getAll`.

Now that our TypeScript logic is done, here is our HTML code:

#### ./ui/src/app/todo/todo.component.html
```html
<h3> Todos </h3>
<table class="table">
  <thead>
    <tr>
      <th>ID</th>
      <th>Description</th>
      <th>Complete</th>
    </tr>
  </thead>
  <tbody>
    <tr *ngFor="let todo of activeTodos">
      <td>{{todo.id}}</td>
      <td>{{todo.message}}</td>
      <td>
        <button *ngIf="!todo.complete" class="btn btn-secondary" (click)="completeTodo(todo)">
          <i class="fa fa-check"></i>
        </button>
        <button *ngIf="todo.complete" class="btn btn-success" disabled>
          <i class="fa fa-check"></i>
        </button>

        <button class="btn btn-danger" (click)="deleteTodo(todo)">
            <i class="fa fa-trash"></i>
        </button>
      </td>
    </tr>
  </tbody>
</table>
<h3>Completed</h3>
<table class="table">
    <thead>
      <tr>
        <th>ID</th>
        <th>Description</th>
        <th>Complete</th>
      </tr>
    </thead>
    <tbody>
      <tr *ngFor="let todo of completedTodos">
        <td>{{todo.id}}</td>
        <td>{{todo.message}}</td>
        <td>
          <button *ngIf="!todo.complete" class="btn btn-secondary" (click)="completeTodo(todo)">
            <i class="fa fa-check"></i>
          </button>
          <button *ngIf="todo.complete" class="btn btn-success" disabled>
            <i class="fa fa-check"></i>
          </button>
  
          <button class="btn btn-danger" (click)="deleteTodo(todo)">
              <i class="fa fa-trash"></i>
          </button>
        </td>
      </tr>
    </tbody>
  </table>
<input placeholder="description..." [(ngModel)]="todoMessage">
<button class="btn btn-primary" (click)="addTodo()"> Add </button>
```

Our todo html is essentially two tables, that use the `*ngFor` directive to iterate over all the todo items in our variable activeTodos and completedTodos, which have been set by the `getAll` function. The table then contains the id and message of the todo item, as well as two button which will complete and delete the given todo. We also indicate whether the todo item has been completed already. We do this using an `*ngIf` directive. If the `todo.complete` property is false, we show an active green button, and if the todo is already completed, we show a grey disabled button.

```
NOTE: It would be best practice to extract these tables as a component by themselves. But for the simplicity of this tutorial, I have chosen to keep this in a single file.
```

At the very bottom of the html page, we have an input string and a button whose click is mapped to our `addTodo` function. This is what gives our users the possibility of adding new todo items. The input content is mapped to our `todoMessage` variable via. our ngModel directive. This directive works like a two-way binding, meaning that the variable is tied to the input element and the element is tied to the variable, should one change, so will the other. This is why our `addTodo` function creates a new Todo item, using the `todoMessage` variable.

If we were to spin up our project now, none of our hard work would show. We aren't routing our clients anyway and even if we were, we are not authenticated to get the todo list from our backend. So, let's setup routing together with an authentication service next, so we can get access to our data.

### Setting up Routing and Securing it with Auth0
We will start by creating our authentication service. This service will be using the Auth0 CDN import, which was inserted at the start of this article. Keep in mind, that this is an in-memory authentication service, the authentication token that we retrieve will not be saved anywhere, so if a user is to reload the web page, they will haev to re-authorize against auth0.

```
NOTE: insert something here about localStorage
```

Either way, our authentication service, will look like this:

#### ./ui/src/app/service/auth.service.ts
```js
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { HttpClient, HttpErrorResponse, HttpHeaders } from "@angular/common/http";
import { Router } from "@angular/router";

import * as auth0 from 'auth0-js';

(window as any).global = window;

@Injectable()
export class AuthService {
    constructor(public router: Router)  {}

    access_token: string;
    id_token: string;
    expires_at: string;

    auth0 = new auth0.WebAuth({
        clientID: 'bpF1FvreQgp1PIaSQm3fpCaI0A3TCz5T',
        domain: environment.domain,
        responseType: 'token id_token',
        audience: 'https://' + environment.domain + '/userinfo',
        redirectUri: environment.callback,
        scope: 'openid'
    });

    public login(): void {
        this.auth0.authorize();
    }

    // ...
    public handleAuthentication(): void {
        this.auth0.parseHash((err, authResult) => {
            if (authResult && authResult.accessToken && authResult.idToken) {
                window.location.hash = '';
                this.setSession(authResult);
                this.router.navigate(['/home']);
            } else if (err) {
                this.router.navigate(['/home']);
                console.log(err);
            }
        });
    }

    private setSession(authResult): void {
        // Set the time that the Access Token will expire at
        const expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
        this.access_token = authResult.accessToken;
        this.id_token = authResult.idToken;
        this.expires_at = expiresAt;
    }

    public logout(): void {
        this.access_token = null;
        this.id_token = null;
        this.expires_at = null;
        // Go back to the home route
        this.router.navigate(['/']);
    }

    public isAuthenticated(): boolean {
        // Check whether the current time is past the
        // Access Token's expiry time
        const expiresAt = JSON.parse(this.expires_at || '{}');
        return new Date().getTime() < expiresAt;
    }

    public createAuthHeaderValue(): string {
        if (this.id_token == "") {
            "";
        }
        return 'Bearer ' + this.id_token;
    }
}
```

As you can see, we have the familiar `@Injectable()` decorator, which we use for injecting an angular Router into our service on initialisation. We then define three string properties and a single `auth0.WebAuth` object. This object is what is used for authenticating against Auth0, so we will need to use the information from our environment files. The `auth0.WebAuth` object will in turn use this information to send to Auth0 and inform which application we are trying to get access to.

Our `login` function is quite simple looking, but essentially, this will initialise authentication with Auth0, redirecting the client to the Auth0 login site. If the user is authenticated, they will be sent to the specified callback location, with a tailing hash on the URL. This url will be parsed by the `handleAuthentication` function, which if successful calls the `setSession` function, which very simple sets our three properties with the appropriate values.

The `lougout` function resets all tokens from memory and the `isAuthenticated` method just returns whether the token as expired or not. We will use this later, for getting an authentication status of our client.

Lastly, we have our `createAuthHeaderValue`, returns a string in the form of an `Authorization` bearer header. In other words, it appends the `id_token` property, to 'Bearer '.

Easy peasy! Just remember to add this service to the `providers` section of our `app.module.ts` file.

Now we need to create our callback component, which will be done like the other components:

> ng g c callback

Change the `callback.component.ts` to the following:

#### ./ui/src/app/callback/callback.component.ts
```js
import { Component, OnInit } from '@angular/core';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-callback',
  templateUrl: './callback.component.html',
  styleUrls: ['./callback.component.css']
})
export class CallbackComponent implements OnInit {

  constructor(private auth: AuthService) { }

  ngOnInit() {
    this.auth.handleAuthentication();
  }
}
```

All we are doing, is invoking the `handleAuthentication` method of our authentication service, once the component is initialised, parsing the url hash and redirecting the client to '/home'.

```
If you want you can add your favourite loading gif to this component, to let your users know that something is happening. However, this is not mandatory.
```

Once we have set our authentication session with our callback component, we need to make sure that all future requests include our retrieved token, in the `Authorization` header. To do this, we will create yet another service. So, create a new file: `./ui/src/app/service/token.interceptor.ts`. This is our interceptor service, which will intercept all of our outgoing requests and add an authentication header, if a token is available.

#### ./ui/src/app/service/token.interceptor.ts
```js
import { Injectable } from "@angular/core";
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent } from "@angular/common/http";
import { AuthService } from "src/app/service/auth.service";
import { Observable } from "rxjs/internal/Observable";

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(public auth: AuthService) {}
  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    request = request.clone({
      setHeaders: {
        Authorization: this.auth.createAuthHeaderValue()
      }
    });
    return next.handle(request);
  }
}
```

This service implements HttpInterceptor, which has the constraint of needing the `intercept` function, which will be revoked upon an http request. Much like our middleware in our backend, this will return a 'next request', which basically forwards the modified request to the original destination. The `intercept` function in this case is quite simple. We are cloning in the incoming request, but adding an Authorization header, using our `auth.createAuthHeaderValue` from our AuthService. Remember, if no token is found, nothing is added. Once the request has been altered, we invoke the `next.handle` passing the modified request.

The reason why we are creating this service, is to ensure that all of our requests contain the appropriate Authorization header. Centralising the header management like this, makes it easier in the future, if changes are made to authentication or if new requests are added to the project.

To include this in our `app.module.ts` file, we will need to add the following object to our providers array, importing all the necessary dependencies:

```js
{
    provide: HTTP_INTERCEPTORS,
    useClass: TokenInterceptor,
    multi: true
}
```

Almost there! Now, we need to change our `app.component.ts` page to include routing and actually define our routes. Let's start by creating our routing definition. Create a new file: `./ui/src/app/app-routing.module.ts` and in this file, write the following code:

#### ./ui/src/app/app-routing.module.ts
```js
import { HomeComponent } from "./home/home.component";
import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import { AuthGuardService } from "./service/auth-guard.service";
import { CallbackComponent } from "./callback/callback.component";
import { TodoComponent } from "./todo/todo.component";

const routes: Routes = [
    { path: '', redirectTo: 'home', pathMatch: 'full' },
    { path: 'home', component: HomeComponent },
    { path: 'todo', component: TodoComponent,  canActivate: [AuthGuardService] },
    { path: 'callback', component: CallbackComponent }
  ];
  
  @NgModule({
    imports: [ RouterModule.forRoot(routes) ],
    exports: [ RouterModule ]
  })
  export class AppRoutingModule { }
```

The only important aspect of this, if our constant routes. This defines an array of paths. Our root path will redirect to our HomeComponent (and os will '/home'). Our todo path will redirect to our TodoComponent and callback to our CallbackComponent. However, you will notice, that the todo path is slightly different, in that it is using the `canActivate` property and using the `AuthGuardService`. Now... we haven't written our AuthGuardService yet, so let's do that right away.

The `canActivate` property basically asks a service (which implements CanActivate) for a boolean response of whether or not a user is able to activate this page. So let's quickly create this service as well. Create a new file:  `./ui/src/app/service/auth-guard.service.ts`, which will contain the logic to whether or not a user is authentorized.

#### ./ui/src/app/service/auth-guard.service.ts
```js
import { CanActivate, Router } from "@angular/router";
import { Injectable } from "@angular/core";
import { AuthService } from "./auth.service";
import { environment } from "src/environments/environment";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs";


@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(private auth: AuthService, private router: Router) {}

    canActivate(): boolean {
        if (this.auth.isAuthenticated()) {
            return true;
        }

        this.auth.login()
    }
}
```

Really a quite simple service, which calls our AuthService, asking whether the current user is authenticated. If the user is authenticated, we will return a true forwarding the user to the request route. If not, we will invoke the `auth.login` function, which will send our user to the Auth0 login page.

Now, we will include our routing and our auth service in our `app.module.ts` file, so the finalised version will look as such:

#### ./ui/src/app/app.module.ts
```js
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { AuthGuardService } from './service/auth-guard.service';
import { AuthService } from 'src/app/service/auth.service';
import { CallbackComponent } from 'src/app/callback/callback.component';
import { TodoComponent } from './todo/todo.component';
import { TodoService } from './service/todo.service';
import { FormsModule } from '@angular/forms';
import { TokenInterceptor } from 'src/app/service/token.interceptor';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    CallbackComponent,
    TodoComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    HttpClientModule
  ],
  providers: [AuthGuardService, AuthService, TodoService, {
    provide: HTTP_INTERCEPTORS,
    useClass: TokenInterceptor,
    multi: true
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
```

For the routing, as you can see, we have imported our `AppRoutingModule`. So, now we can finally include this routing to our application, which will be the last part of this tutorial.

## Putting it all together
So, it's been a long while, but now we are finally going to add routing to the application. As well as creating a little menu bar, for navigation and logging out. First we need to edit the `./ui/src/app/app.component.ts` to add our AuthService to our AppComponent.

#### ./ui/src/app/app.component.ts
```js
import { Component } from '@angular/core';
import { AuthService } from 'src/app/service/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {

  constructor(public auth: AuthService) {}
}
```

Next, we need to edit the `./ui/src/app/app.component.html` file and set it to the following code:

#### ./ui/src/app/app.component.html
```html
<nav class="navbar navbar-default">
    <div class="container-fluid">
      <div class="navbar-header">
        <a class="navbar-brand" href="#">Auth0 - Angular</a>
  
        <button
          class="btn btn-light btn-margin"
          routerLink="/">
            Home
        </button>
  
        <button
          class="btn btn-light btn-margin"
          *ngIf="!auth.isAuthenticated()"
          (click)="auth.login()">
            Log In
        </button>
  
        <button
          class="btn btn-light btn-margin "
          *ngIf="auth.isAuthenticated()"
          (click)="auth.logout()">
            Log Out
        </button>
      </div>
    </div>
  </nav>
  
<router-outlet></router-outlet>
```

We are creating a new nav bar including functions from our auth service to login and logout. Underneath this nav bar, we are including the `router-outlet` component. This is what tells Angular to ask the routing module, which page it should load.

```
NOTE: As with our todo list tables, following best practice, we would have created a new component for our navigation bar.
```

It's a very common practice to include the navigation bar together with the routing component. This ensures that the menu is present on all of our pages, without having to reload the navigation bar for each individual page.

However, we are now done! Our todo list is working as intended and we can now rejoice in creating a todo list that only users which we have specified are granted access to. The last thing we need to do is build our ui and start our web server. Build the UI with the following command:

> ng build --prod 

This will place a few transpiled and compressed javascript files in the `./ui/dist` folder. Which (luckily) is where we are serving our static files from, with our web server.

So, go to the root of the project and use the command:

> go run main.go

And the server will be up and running! Go to `localhost:3000/` and celebrate by creating all your secured todo items!

```
NOTE: You can also compile the go code into a binary using the command: go build -o todo-server main.go, which will create an executable file called todo-server. The file will be compiled to your system, but that can be modified using the GOOS and GOARCH environment variables, read more here: https://golang.org/cmd/compile/
```

## Conclusion
So, we finally made it! The application itself that we created was pretty simple. Just a todo list, where we can add delete and complete some todo items. However, the framework around our application is quite sound. We have handled authentication via. Auth0, creating a very strong starting point for our application, by starting with security in mind.

Adding features to our application becomes a lot easier, once we have established a strong fundament in security. We can add different todo lists for different users, relatively easily, without having to worry about how this will affect our application down the road. Using a third party security solution like Auth0, is also a great advantage, because we can rest assured that this solution will keep our application data safe. With a few changes here and there (such as serving our API and static files over HTTPS), we could quite confidently deploy this code to production. 

I hope this article has been helpful, and has given some insight to how easy it is to implement Auth0 as a third-party authentication service, as well as using Angular as a frontend and Golang as a backend. Feedback and questions are very welcome!
