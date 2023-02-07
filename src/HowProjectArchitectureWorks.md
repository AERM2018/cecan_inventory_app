# Project's architecture (directories distribution)

This project implements an architecture (*arch*) model called "Clean Architecture" which divides the project into three different layers. Thanks to this arch it's easier
to separate the project's parts by their functionality. The layers are the following:
- Infrastructure
- Adapters
- Domain


## Infrastructure

This layer contains the configuration and initialization of the project and it's dependencies (third-party software). The dependencies are the software that I used to be build this software.
Within this folder there are others folder such as config, external, http and storage.

---------------------------
#### config
This folder contains the file of the app´s configuration where the application instance is initialized with the routes and middlewares, also has the method which starts the server.

#### external
This folder contains different function that are part of the third-party software that was used. Within this you will find:
- AuthToken : Function to generate the JWT (Json Web Token). Token that is used to authenticate the user in the server-side, by doing this it's possible to know which users are allowed to do any action defined in the API.
- BodyReader: Function to read the request's body and pass it to the api endpoints' middlewares until get the request's handler.
- DataSources: Objects that implement functions that are able to read, create, update and delete records from the data base (most of the cases, one data source is for a table in the data base)


#### http
This folder contains the files that define the routes the API is listening to, as well as,  the middlewares that are functions that executed before execute the main handler. It looks like this:
- Middlewares
  - requestValidator: Function to validate a request contains the specified fields
  - customReqValidator: Contains different functions that specify the fields that has to be in the reques's body and the validations that they have to fulfill
  - verifyJwt: Function to validate the JWT token that is present in the request's headers, if it's valid, the next middlewares and main handler are executed. If not, a message that the token is invalid is sent to the client.
  - DbValidator: Object that make calls to the data base to check if certain data is in the data base since it's needed for executing the different actions that the api is able to
- Routes
  This folder contains the different files in which the routes are difined, each file groups certain routes that belong to a module, for example, the pharmacy's routes are for getting the catalog and inventory from that department.
  
 #### storage
 This folder contains different files that are useful to set up the data base connection. Within it you will find:
 - Migrations: Folder in which all the migration files are located. The migration files are SQL files that have SQL commands to create, alter or delete the data base tables. Thanks to this you can run the API the first time, and the data base tables will be created automatically.
 - Migrator: Object that has the different function to run or destroy the data base migrations.
 - Seeds: They are kind of migration that insted of creating a table, it will insert data to a table when the app runs.
 - Storage: This is the file where the connection is stablished and the migrations run to read and execute the missing migrations and seeds. It uses the environmental variables to define the data base credentials to be used.

------------------------------

## Adapters
This folder contains two subfolders that are shown below:

#### helpers
This folder contains two function that are shown below:
- responseHandler: A function that receives a group of parameters that are part of the response that is sent to the client, this is the function in charge of triggering the event to send the response.
- uploadFile: This a function to take the file that is present in the request's body, stores it and return the path where it was stored. If there's an error, it will call the another helper to send an error to the client.

#### controllers
This folder contains different files that are the main handlers of the routes, they are many files, but all of them follows the same logic:
- The handler is an object that has as parameters the data sources instances and interactors that were injected when the controller instance was defined in the router
  - You can the interacors explanations here:
- It calls the desired method of the interactor, stores the response from it that is actually the response to the client, and calls the helper to send the response

--------------------------------

## Domain

#### assets
This folder contains assets (images) that are used in the creation of the different documents.

#### common
This folder contains functions that are re-used in different parts of the project such as 
- get the semaforization color by date.
- insert a predefined header to pfd documents.
- mailer service.
- upload data from files to data base

#### mocks
This folder containt objects that contain objects with the data thar is inserted to the data base the first time is run.

#### models
This folder contains objects that are a representation of the data bases tables in the programming language's objects. These objects have also a function that are executed before inserting or deleteing records in the corresponding table.

#### pfds
This folder is used to store the pdf files that are created whent the user requests the generation of a document. Those files can be deleted each lapse of time in order not to make this folder heavy.

#### seeds
This folder is used to store the files that the user upload to insert predifend data to a table (until now it's possible to upload data to fixed assests table).

#### useCases
This folder contains different files that contain the business login that is the responsible of dictate the way the data must be treated, they are many files, but all of them follows the same logic:
- It's an object call interactor, its intance is declared in the controller and the data sources are injected to it.
- It contains functions that are called in the controllers, each function makes a call to the data base using a data source.
- Thanks to the execution of the login within it, it's possible to determine the response that will be sent to the client by the controller. In other words, the response body is returned by this function
