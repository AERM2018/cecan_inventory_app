
# Table of content
- [Introducción](#introduction)
- [Infrastructure](#infrastructure)
  - [config](#sub_config)
  - [external](#sub_external)
  - [http](#sub_http)
  - [storage](#sub_storage)
- [Adpaters](#adapters)
  - [helpers](#sub_helpers)
  - [controllers](#sub_controllers)
- [Domain](#domain)
  - [assets](#sub_assets)
  - [common](#sub_common)
  - [mocks](#sub_mocks)
  - [models](#sub_models)
  - [pfds](#sub_pfds)
  - [seeds](#sub_seeds)
  - [useCases](#sub_useCases)

<div id='introduction'/>

# Arquitectura del proyecto (distribución de directorios)

Este proyecto implementa un modelo de arquitectura (*arch*) llamado "Arquitectura Limpia" que divide el proyecto en tres capas diferentes. Gracias a este arquitectura es más fácil separar las partes del proyecto por su funcionalidad. Las capas son las siguientes
- Infrastructure
- Adapters
- Domain


<div id='infrastructure'/>

## Infrastructure

Esta capa contiene la configuración e inicialización del proyecto y sus dependencias (software de terceros). Las dependencias son el software que he utilizado para construir este software.
Dentro de esta carpeta hay otras carpetas como config, external, http y storage.

---------------------------

<div id='sub_config'/>

### config
Esta carpeta contiene el archivo de configuración de la aplicación donde se inicializa la instancia de la aplicación con las rutas y middlewares, también tiene el método que inicia el servidor.

<div id='sub_external'/>

### external
Esta carpeta contiene diferentes funciones que son parte del software de terceros que se utilizó. Dentro de esta se encuentran:
- AuthToken : Función para generar el JWT (Json Web Token). Token que se utiliza para autenticar al usuario en el lado del servidor, haciendo esto es posible saber que usuarios tienen permiso para realizar cualquier acción definida en el API.
- BodyReader: Función para leer el cuerpo de la petición y pasarlo a los middlewares de los endpoints de la api hasta obtener el manejador de la petición.
- DataSources: Objetos que implementan funciones capaces de leer, crear, actualizar y borrar registros de la base de datos (en la mayoría de los casos, un data source es para una tabla de la base de datos).


<div id='sub_http'/>

### http
Esta carpeta contiene los archivos que definen las rutas que la API está escuchando, así como, los middlewares que son funciones que se ejecutan antes de ejecutar el manejador principal. Su aspecto es el siguiente
- Middlewares
  - requestValidator: Función para validar que una petición contiene los campos especificados.
  - customReqValidator: Contiene diferentes funciones que especifican los campos que tienen que estar en el cuerpo de la petición y las validaciones que tienen que cumplir
  - verifyJwt: Función para validar el token JWT que está presente en las cabeceras de la petición, si es válido, se ejecutan los siguientes middlewares y el handler principal. En caso contrario, se envía un mensaje al cliente indicando que el token no es válido.
  - DbValidator: Objeto que realiza llamadas a la base de datos para comprobar si ciertos datos están en la base de datos ya que son necesarios para ejecutar las diferentes acciones que es capaz de realizar la api
- Routes
  Esta carpeta contiene los diferentes archivos en los que se definen las rutas, cada archivo agrupa ciertas rutas que pertenecen a un módulo, por ejemplo, las rutas de farmacia son para obtener el catálogo e inventario de ese departamento.
  
  <div id='sub_storage'/>

### storage
 Esta carpeta contiene diferentes archivos útiles para configurar la conexión a la base de datos. Dentro de ella encontrará:
 - Migrations: Carpeta en la que se encuentran todos los ficheros de migración. Los ficheros de migración son ficheros SQL que contienen comandos SQL para crear, alterar o borrar las tablas de la base de datos. Gracias a esto podrás ejecutar la API la primera vez, y las tablas de la base de datos se crearán automáticamente.
 - Migrator: Objeto que tiene funciones diferentes para ejecutar o destruir las migraciones de base de datos.
 - Seeds: Son un tipo de migración que en lugar de crear una tabla, insertará datos a una tabla cuando se ejecute la aplicación.
 - Storage: Es el fichero donde se establece la conexión y se ejecutan las migraciones para leer y ejecutar las migraciones y semillas que faltan. Utiliza las variables de entorno para definir las credenciales de la base de datos a utilizar.

------------------------------

<div id='adapters'/>

## Adapters
Esta carpeta contiene dos subcarpetas que se muestran a continuación:

<div id='sub_helpers'/>

### helpers
- responseHandler: Es una función que recibe un grupo de parámetros que forman parte de la respuesta que se envía al cliente, esta es la función encargada de disparar el evento para enviar la respuesta.
- uploadFile: Es una función que toma el archivo presente en el cuerpo de la petición, lo almacena y devuelve la ruta donde fue almacenado. Si hay un error, llamará al otro helper para enviar un error al cliente.

<div id='sub_controllers'/>

### controllers
Esta carpeta contiene diferentes archivos que son los manejadores principales de las rutas, son muchos archivos, pero todos siguen la misma lógica:
- El manejador es un objeto que tiene como parámetros las instancias de los data sources y los interactors que se inyectaron cuando se definió la instancia del controlador en el router.
  - Puedes ver las explicaciones de los interactores aquí:
- Llama al método deseado del interactor, almacena la respuesta del mismo que es en realidad la respuesta al cliente, y llama al helper para enviar la respuesta

--------------------------------

<div id='domain'/>

## Domain

<div id='sub_assets'/>

### assets
Esta carpeta contiene activos (imágenes) que se utilizan en la creación de los diferentes documentos.

<div id='sub_common'/>

### common
Esta carpeta contiene funciones que se reutilizan en diferentes partes del proyecto como por ejemplo 
- obtener el color de semaforización por fecha.
- insertar una cabecera predefinida a los documentos pfd.
- servicio de correo.
- subir datos de ficheros a base de datos

<div id='sub_mocks'/>

### mocks
Esta carpeta contiene objetos con los datos que se insertan en la base de datos la primera vez que se ejecuta.

<div id='sub_models'/>

### models
Esta carpeta contiene objetos que son una representación de las tablas de las bases de datos en los objetos del lenguaje de programación. Estos objetos tienen también una función que se ejecuta antes de insertar o borrar registros en la tabla correspondiente.

<div id='sub_pdfs'/>

### pfds
Esta carpeta se utiliza para almacenar los ficheros pdf que se crean cuando el usuario solicita la generación de un documento. Estos ficheros pueden ser borrados cada cierto tiempo para no hacer pesada esta carpeta.

<div id='sub_seeds'/>

### seeds
Esta carpeta se utiliza para almacenar los ficheros que el usuario sube para insertar datos predefinidos en una tabla (hasta ahora es posible subir datos a tablas de activos fijos).

<div id='sub_useCases'/>

### useCases
Esta carpeta contiene diferentes archivos que contienen la logíca de negocio que es el encargado de dictar la forma en que los datos deben ser tratados, son muchos archivos, pero todos siguen la misma lógica:
- Es un llamado objeto interactor, su instancia se declara en el controlador y se le inyectan los data sources.
- Contiene funciones que son llamadas en los controladores, cada función hace una llamada a la base de datos utilizando un data source.
- Gracias a la ejecución de la logíca dentro de él, es posible determinar la respuesta que será enviada al cliente por el controlador. En otras palabras, el cuerpo de la respuesta es devuelto por esta función