# CECAN Inventory app API
This is the API REST that the CECAN Inventory app uses to manage the data stored in a data base. It written in the Go programming Language

## How to run it
This server-side app is executed in docker container, so it's needed to have docker install already in order to be able to run it.


1. Being in the project's folder root go to the directory Docker/dev
2. Create the files .env and .postgres.env with the fields that you can find within the examples in the root folder with the same name plus a .example word (.env.example & .env.postgres.example). Those files contain the environmental variables that are used to stablish the data base connection and other important fields
3. Once the variables are defined in the files mentioned before, you have to run the next command
```
  docker-compose -p cecan-dev up -d
```
4. After some seconds you should see a message that says the containers have been created successfully.
5. Run the command: ```docker ps```, the you will see the list of the containers that are currently running, you should see two containers which their name starts with "cecan-dev"
6. Finally, you can access to the container on the ports specified, 
in the table that displays the command ```dokcer ps``` you will see something like ```0.0.0.0:4000-->4000```, 
the left side of the string corresponds to the port which you use on the local machine to connect to the container, and the right side is the port inside the container that is being used.
