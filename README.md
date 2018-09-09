
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

# Running the application
You can either run the application using an in-memory database, or using a postgreSQL database. To use the application with the in-memory database, use the following command:

> go run main.go -audience=<auth0_audience> -domain=<auth0_domain> -db=memory

To use it with the postgresql, you need a configuration file, with the following content:

#### .\config.json
```
{
    "host": "172.16.1.68",
    "port": 5432,
    "user": "postgres",
    "password": "postgres",
    "db_name": "todo" 
}
```

...and then run the application with the following command:

> go run main.go -audience=<auth0_audience> -domain=<auth0_domain> -db=postgres -config="config.json"