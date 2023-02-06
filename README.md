# API de la aplicación Inventario CECAN
Esta es la API REST que la app Inventario CECAN utiliza para gestionar los datos almacenados en una base de datos. Está escrita en el lenguaje de programación Go

## Cómo ejecutarla
Esta aplicación de servidor se ejecuta en un contenedor docker, por lo que es necesario tener docker instalado para poder ejecutarla.


1. Estando en la raíz de la carpeta del proyecto ve al directorio Docker/dev
2. Crea los ficheros .env y .postgres.env con los campos que puedes encontrar dentro de los ejemplos en la carpeta raíz con el mismo nombre más una palabra .example (.env.example & .env.postgres.example). Esos archivos contienen las variables de entorno que se usan para establecer la conexión a la base de datos y otros campos importantes
3. Una vez definidas las variables en los ficheros antes mencionados, hay que ejecutar el siguiente comando
```
  docker-compose -p cecan-dev up -d
```
4. Después de unos segundos deberías ver un mensaje que dice que los contenedores se han creado correctamente.
5. Ejecuta el comando ```docker ps```, el verás la lista de los contenedores que se están ejecutando actualmente, deberías ver dos contenedores cuyo nombre empieza por "cecan-dev"
6. Finalmente, puedes acceder al contenedor en los puertos especificados, 
en la tabla que muestra el comando ```dokcer ps``` verás algo como ```0.0.0.0:4000-->4000```, 
el lado izquierdo de la cadena corresponde al puerto que usas en la máquina local para conectarte al contenedor, y el lado derecho es el puerto dentro del contenedor que se está usando.
