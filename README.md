# go-mongo-restapi
ApiREST making use of own [scv-go-framework](https://github.com/scanet9/scv-go-framework).

Provides:
- Basic CRUD functionalities for user management
- MongoDB persistent storage
- JWT token-based authorization.
## To run it locally:
```
    go run cmd/main.go
```
or:
```
cd cmd
go build main.go
./main
```

## To run it in a docker container:
```
docker build -t go-mongo-restapi .
docker run --name go_mongo -p 8080:8080 go-mongo-restapi
docker stop go_mongo
```