# go-mongo-restapi
REST API making use of own [scv-go-framework](https://github.com/sergicanet9/scv-go-framework) package.

Provides:
- Basic CRUD functionalities for user management
- MongoDB persistent storage
- Multi-environment configs
- JWT token-based authorization
- Swagger UI documentation
- Dockerized app

## Run the application locally
```
    go run cmd/main.go -env={env}
```
or:
```
go build cmd/main.go
 ./main -env={env}
```
Then open http://localhost:{port}/swagger/index.html in the browser.
<br />
<br />
 NOTES:
- The env flag is optional. Default value: "local".
- It is also possible to debug it in Visual Studio Code with the provided [launch.json](https://github.com/sergicanet9/go-mongo-restapi/blob/main/.vscode/launch.json).

## Run the application in a docker container
```
docker build -t go-mongo-restapi .
docker run --name {container_name} -p {port}:{port} -e env={env} go-mongo-restapi
```
Then open {address}:{port}/swagger/index.html in the browser.
<br />
<br />
NOTES:
- The env flag is optional. Default value: "docker".
- For running the image in an Azure Web App in port 443 instead of 80 (only these two ports are exposed), it´s necessary to add an entry on the Web App´s application settings with the port´s value (443) and name it "WEBSITES_PORT". Also make sure to change the value in [config.docker.json](https://github.com/sergicanet9/go-mongo-restapi/blob/main/config/config.docker.json).

### To stop the running docker container
```
docker stop {container_name}
```
## (Re)Generate Swagger documentation
```
swag init -g cmd/main.go
```
