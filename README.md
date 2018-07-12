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