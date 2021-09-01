# go-mongo-restapi
REST API making use of own [scv-go-framework](https://github.com/scanet9/scv-go-framework).

Provides:
- Basic CRUD functionalities for user management
- MongoDB persistent storage
- JWT token-based authorization.
- Swagger UI documentation.

## How to run it locally:
## 1. Regenerate Swagger (only in case of any code change)
```
swag init -g cmd/main.go
```
## 2. Run or build the application

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

## How to run it in a docker container:
```
docker build -t go-mongo-restapi .
docker run --name go_mongo -p 8080:8080 go-mongo-restapi
```
Then open http://localhost:8080/swagger/index.html in the browser.
NOTE: There is no need to manually generate the swagger documentation, the [Dockerfile](https://github.com/scanet9/go-mongo-restapi/blob/main/Dockerfile) will do it.

### To stop the docker container:
```
docker stop go_mongo
```