# go-mongo-restapi
REST API making use of own [scv-go-framework](https://github.com/scanet9/scv-go-framework).

Provides:
- Basic CRUD functionalities for user management
- MongoDB persistent storage
- JWT token-based authorization.
- Swagger UI documentation.

## Run the application locally
```
    go run cmd/main.go
```
or:
```
cd cmd
go build main.go
./main
```
Then open http://localhost:8080/swagger/index.html in the browser.
<br />
NOTE: It is also possible to debug it in Visual Studio Code with the provided [launch.json](https://github.com/scanet9/go-mongo-restapi/blob/main/.vscode/launch.json).

## Run the application in a docker container
```
docker build -t go-mongo-restapi .
docker run --name go_mongo -p 8080:8080 go-mongo-restapi
```
Then open http://localhost:8080/swagger/index.html in the browser.
### To stop the docker container
```
docker stop go_mongo
```
## (Re)Generate Swagger documentation
```
swag init -g cmd/main.go
```
