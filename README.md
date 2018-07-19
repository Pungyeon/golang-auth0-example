# Creating a Todo list with Authentication, using Golang, Angular and Auth0

## What we will be building

## Prequisites & Installation

## The Todo list backend
### Creating our in-memory todo list
### Initilialising our Web Server & Serving static files
### Writing the API Endpoints
### Using middleware for authorization

## Creating our Frontend
### Initialising our angular project
### Writing a welcome page
### Writing our todo list functionality
### Creating our auth(0) service
### implementing our auth-guard

## Putting it all together

-- END OF OUTLINE --

# Prerequisites
Please make sure to have nodejs, npm and golang installed. The angular cli is also an advantage, though not necessarily a prerequisite. 
To install the angular cli run the following npm command:

> npm install -g @angular/cli

# (Re)Building the UI
Go to the `ui` directory and run

> npm install

Once npm has finished installing the entirety of the internet, run:

>  ng build --prod

This will place all minimised javascript files in the `dist` folder of the root directory

# Runnign the application
You can run the application using the command

> run main.go

!!WARNING - as of right now, the UI cannot be run with `ng serve`, this will be fixed asap.