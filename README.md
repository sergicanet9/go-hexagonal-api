# go-hexagonal-api
REST API implementing Hexagonal Architecture (Ports & Adapters) making use of own [scv-go-tools](https://github.com/sergicanet9/scv-go-tools) package.

Provides:
- MongoDB and PostgreSQL repository adapters for persistent storage decoupled from the business logic
- Database migrations with Goose for PostgreSQL implementation
- CRUD functionalities for user management
- Multi-environment configs
- JWT authentication and claim-based authorization
- Swagger UI documentation
- Dockerized app
- Integration tests
- CI/CD with Github Actions
- Async process for periodical health checking

## Run and debug the application locally
```
    go run cmd/main.go -v={version} -env={env} -p={port} -db={database}
```
or:
```
go build cmd/main.go
 ./main -env={env} -v={version} -env={env} -p={port} -db={database}
```
Then open {address}:{port}/swagger/index.html in the browser, where {address} is the value specified in the corresponding config.{env}.json, corresponding to "MongoAddress" or "PostgresAddress", depending on the database choosed.
<br />
<br />
 NOTES:
- The v, env, p and db flags are optional. Default values, respectively: "debug", "local", "8080", "mongo".
- For debugging the application with Visual Studio Code´s build-in debugger, select the desired database configuration in the Run and Debug menu and click on Start Debugging button or press F5. All flags can be changed in the debugging configurations provided in [launch.json](https://github.com/sergicanet9/go-hexagonal-api/blob/main/.vscode/launch.json).

## Run the application in a local docker container
```
make mongo-up
```
or:
```
make mongo-up
```
Then open {address}:{port}/swagger/index.html in the browser, where {address} is the value specified in [config.local.json](https://github.com/sergicanet9/go-hexagonal-api/blob/main/config/config.local.json), corresponding to "MongoAddress" or "PostgresAddress", depending on the command run.
<br />
<br />
NOTES:
- When running this commands, the docker images will be built using the current git branch name for the v flag, "local" for the env flag and "8080" for the port flag.

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
make swagger
```

## Database commands for Postgres
### Create a new migration
```
make goose
```
Write the file name without ".sql" suffix and press enter.
Then edit the newly created file to define the behavior of the migration.

### Connect to the remote database
```
psql -h {host} -U {username} {db_name}
```

### Dump the remote database
```
pg_dump -h {host} -U {username} {db_name} --schema-only > dump.sql
```

### Drop the remote database (Azure Postgres Flexible Server)
```
az login
az postgres flexible-server db delete -g {resource_group} -s {resource_name} --database-name {db_name}
```

## Deploy a new environment in Azure
The steps for creating and deploying a new cloud environment on Azure are the following:
1. Create an Azure Web App on Azure Portal and name it go-hexagonal-api-{db}-{env}, where {db} and {env} are the database used and the name of the new environment, respectively.
2. Add an entry on the Web App´s Configuration with name "WEBSITES_PORT" and value "443".
3. On the Web App´s App Service Logs, select "File system" and configure the Quota (MB) and the Retention Period (Days) for enabling app logs.
4. On the Web App´s Deployment Center, select Github Actions option. Authorize the Github account if required. Then select the repository name, the registry source (Azure Container Registry, Docker Hub, etc.) and type the image name in the following format: go-hexagonal-api-{db}.{env}. Make sure that the secrets for accessing to the registry source are provided on Github´s repository settings.
5. Download the publish profile on Web App´s overview and add a secret entry on Github´s repository settings with its value and name it "AZUREAPPSERVICE_PUBLISHPROFILE_{DB}_{ENV}".
6. On Github´s repository settings, add a new environment named {db}-{env} and configure a manual approval protection rule.
7. On the source code, create a file named config.{env}.json in the config folder and add the environment values in it.
8. Edit the CI/CD Pipeline file to include build-{db}-{env} and deploy-{db}-{env} jobs, using the proper environment names and secrets. Make sure that the running port specified in the file is "443".

## Author
Sergi Canet Vela

## License
This project is licensed under the terms of the MIT license.