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
 
