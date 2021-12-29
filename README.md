# go-mongo-restapi
REST API making use of own [scv-go-framework](https://github.com/sergicanet9/scv-go-framework) package.

Provides:
- CRUD functionalities for user management
- MongoDB persistent storage
- Multi-environment configs
- JWT authentication and claim-based authorization
- Swagger UI documentation
- Dockerized app
- Integration tests

## Run the application locally
```
    go run cmd/main.go -env={env}
```
or:
```
go build cmd/main.go
 ./main -env={env}
```
Then open {address}:{port}/swagger/index.html in the browser, where {address} and {port} are the values specified in the corresponding config.{env}.json.
<br />
<br />
 NOTES:
- The env flag is optional. Default value: "local".
- It is also possible to debug it in Visual Studio Code with the provided [launch.json](https://github.com/sergicanet9/go-mongo-restapi/blob/main/.vscode/launch.json), in which the env flag can be modified as well.

## Run the integration tests
```
go test ./test -coverpkg=./... -coverprofile=test/coverage.out
```
Then to see the coverage report run the following command:
```
go tool cover -html=test/coverage.out
```
 NOTES:
- The docker daemon needs to be up and running for executing the tests.

## Run the application in a docker container
```
docker-compose up
```
Then open {address}:{port}/swagger/index.html in the browser, where {address} and {port} are the values specified in the corresponding config.{env}.json.
<br />
<br />
NOTES:
- The env and port values are specified in the [docker-compose.yml](https://github.com/sergicanet9/go-mongo-restapi/blob/main/docker-compose.yml). The port has to match with the one specified in the corresponding config.{env}.json.
- For running the image in an Azure Web App, create a new config.{env}.json and then build the image specifing the new env and port values in the [docker-compose.yml](https://github.com/sergicanet9/go-mongo-restapi/blob/main/docker-compose.yml). Also it is needed to add an entry on the Web App´s application settings with the port´s value and name it "WEBSITES_PORT".

## (Re)Generate Swagger documentation
```
swag init -g cmd/main.go
```

## Author
Sergi Canet Vela

## License
This project is licensed under the terms of the MIT license.
