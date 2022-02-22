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
- CI/CD with Github Actions
- Async process for health checking

## Run and debug the application locally
```
    go run cmd/main.go -env={env} -v={version}
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
- The env and v flags are optional. Default values, respectively: "local", "debug".
- For debugging the application with Visual Studio Code´s build-in debugger, select Run and Debug on the Debug start view or press F5. The env flag can be changed in the debugging configuration provided in [launch.json](https://github.com/sergicanet9/go-mongo-restapi/blob/main/.vscode/launch.json).

## Run the application in a local docker container
```
make up
```
Then open {address}:{port}/swagger/index.html in the browser, where {address} and {port} are the values specified in [config.local.json](https://github.com/sergicanet9/go-mongo-restapi/blob/main/config/config.local.json)
<br />
<br />
NOTES:
- When running this command, the docker image will be built using the value "local" for the env flag and the current git branch name for the v flag.

### Stop and remove the running container
```
make down
```

## Run the integration tests
```
make test
```
Then to see the coverage report run the following command:
```
make cover
```
 NOTES:
- The docker daemon needs to be up and running for executing the tests.

## (Re)Generate Swagger documentation
```
make docs
```

## Deploy a new environment
The steps for creating and deploying a new cloud environment on Azure are the following:
1. Create an Azure Web App on Azure Portal and name it go-mongo-restapi-{env}, where {env} is the name of the new environment.
2. Add an entry on the Web App´s Configuration with name "WEBSITES_PORT" and value "443".
3. On the Web App´s App Service Logs, select "File system" and configure the Quota (MB) and the Retention Period (Days) for enabling app logs.
4. On the Web App´s Deployment Center, select Github Actions option. Authorize the Github account if required. Then select the repository name, the registry source (Azure Container Registry, Docker Hub, etc.) and type the image name in the following format: go-mongo-restapi_api.{env}. Make sure that the secrets for accessing to the registry source are provided on Github´s repository settings.
5. Download the publish profile on Web App´s overview and add a secret entry on Github´s repository settings with its value and name it "AZUREAPPSERVICE_PUBLISHPROFILE_{ENV}".
6. On Github´s repository settings, add a new environment named {env} and configure a manual approval protection rule.
7. On the source code, create two files named docker-compose.{env}.yml and a config.{env}.json in the docker and config folders, respectively, and add the environment values in them. Make sure that the running port is "443" in both files.
8. Modify the CI/CD Pipeline to include build-{env} and deploy-{env} jobs, using the proper environment names and secrets.

## Author
Sergi Canet Vela

## License
This project is licensed under the terms of the MIT license.