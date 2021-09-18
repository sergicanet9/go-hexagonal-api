# go-mongo-restapi
REST API making use of own [scv-go-framework](https://github.com/scanet9/scv-go-framework).

Provides:
- Basic CRUD functionalities for user management
- MongoDB persistent storage
- Multi environment configs
- JWT token-based authorization.
- Swagger UI documentation.

## Run the application locally
```
    go run cmd/main.go -env=local
```
or:
```
go build cmd/main.go
 ./main -env=local
```
Then open http://localhost:8080/swagger/index.html in the browser.
<br />
<br />
 NOTES:
 - The env flag is optional, the default value is "local"
 - It is also possible to debug it in Visual Studio Code with the provided [launch.json](https://github.com/scanet9/go-mongo-restapi/blob/main/.vscode/launch.json).

## Run the application in a docker container
```
docker build -t go-mongo-restapi .
docker run --name go_mongo -p 8082:8082 go-mongo-restapi
```
Then open http://localhost:8082/swagger/index.html in the browser.
### To stop the docker container
```
docker stop go_mongo
```
## (Re)Generate Swagger documentation
```
swag init -g cmd/main.go
```
