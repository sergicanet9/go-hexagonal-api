# go-hexagonal-api
![CI](https://github.com/sergicanet9/go-hexagonal-api/actions/workflows/ci.yml/badge.svg)
![CD](https://github.com/sergicanet9/go-hexagonal-api/actions/workflows/cd.yml/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-80.8%25-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergicanet9/go-hexagonal-api.svg)](https://pkg.go.dev/github.com/sergicanet9/go-hexagonal-api)

A robust REST API for user management built with **Go** and implementing the **Hexagonal Architecture** (Ports & Adapters) pattern. It makes use of own [scv-go-tools](https://github.com/sergicanet9/scv-go-tools) package.
<br />
The API is designed to work seamlessly with either a MongoDB or PostgreSQL database instance, using the same business logic and handlers.

## 🚀 Features
- **Hexagonal Architecture**: Clear separation of concerns with domain, application and infrastructure layers.
- **Database Agnostic**: Decoupled repository adapters allow switching between MongoDB and PostgreSQL storage without changing core logic.
- **Authentication & Authorization**: Implements JWT authentication and claim-based authorization for secure endpoints for user management.
- **Asyncronous Process**: Go routines management with an included process for periodically health checking the application.
- **Testing**: Comprehensive unit tests with code coverage and integration tests for the happy path.
- **Developer Experience**: Built-in Makefile, Swagger UI for API documentation and management UIs for each database (pgAdmin and mongo-express).
- **Lifecycle Management**: Multi-environment support with config files, dockerfile and docker-compose, CI/CD pipelines with GitHub Actions, Kubernetes deployment and New Relic observability.

## 🏁 Getting Started
### Run it with Docker
To start the entire application stack using Docker Compose, run:
```
make up
```
This command launches six containers:
* Two databases (MongoDB and PostgreSQL).
* Two database UIs (mongo-express and pgAdmin).
* Two API instances, one for each database.

The URLs for Swagger and the database UIs will be printed in the console.

To stop and remove all containers, run:
```
make down
```

### Run it with command line
Run a single API instance with command-line arguments with the following command:
```
    go run cmd/main.go --ver={version} --env={environment} --port={port} --db={database} --dsn={dsn} --nrkey={newrelic_key}
```
or:
```
go build cmd/main.go
 ./main --ver={version} --env={environment} --port={port} --db={database} --dsn={dsn} --nrkey={newrelic_key}
```
Provide the desired values for: `{version}`, `{environment}`, `{port}`, `{database}`, `{dsn}`.
<br />
The `--nrkey` flag is optional and can be omitted if you do not want to configure New Relic observability.
<br />
Then open `http://localhost:{port}/swagger/index.html` to access the Swagger UI page.
<br />
<br />
NOTES:
- The target database container needs to be up and running (run `make up`).

### Debug it with VS Code
The project includes debugging profiles in [launch.json](https://github.com/sergicanet9/go-hexagonal-api/blob/main/.vscode/launch.json) for both MongoDB and PostgreSQL setups. Simply select the desired configuration in the VS Code debugger and run it.
<br />
Then open `http://localhost:{port}/swagger/index.html` to access the Swagger UI page.
<br />
<br />
NOTES:
- The target database container needs to be up and running (run `make up`).

## 📦 API Endpoints
### Public Routes
These endpoints don't require authentication.
<br />
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Performs a health check of the API's status. |
| `POST` | `/v1/users` | Creates a new user. |
| `POST` | `/v1/users/login` | Authenticates a user and returns a JWT token. |
| `POST` | `/v1/users/many` | Creates multiple users from a list. |

### Protected Routes
These endpoints require a valid JWT in the Authorization header, formatted as `Bearer {token}`.
<br />
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/v1/users` | Retrieves all users. |
| `GET` | `/v1/users/email/{email}` | Retrieves a user by their email address. |
| `GET` | `/v1/users/{id}` | Retrieves a user by their unique ID. |
| `PATCH` | `/v1/users/{id}` | Updates a user's information. |
| `GET` | `/v1/claims` | Returns all existing claims. |


### Admin Routes
These endpoints require a valid JWT in the Authorization header with the `admin:true' claim.
<br />
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `DELETE` | `/v1/users/{id}` | Deletes a user by their unique ID. |

## ✅ Testing
### Run unit tests with code coverage
```
make test-unit
```

### View coverage report
```
make cover
```

### Run integration tests
```
make test-integration
```
 NOTES:
- Docker is required for running integration tests.

## 🛠️ Other Commands 
### (Re)Generate Swagger documentation
```
make swagger
```
### (Re)Generate Mockery mocks
```
make mocks
```

### Create a new PostgreSQL migration
```
make goose
```
Write the file name without ".sql" suffix and press enter.
Then edit the newly created file to define the behavior of the migration.

### Connect to pgAdmin
Open the pgAdmin URL printed after running `make up`.
<br />
Log in with the email and password specified as `PGADMIN_LOGIN_EMAIL` and `PGADMIN_LOGIN_PASSOWRD` in the [.env](https://github.com/sergicanet9/go-hexagonal-api/blob/main/.env) file.
<br />
When prompted for the database user password, use the value of `POSTGRES_PASSWORD` from the same file.

### Connect to mongo-express
Open the mongo-express URL printed after running `make up`.
<br />
Log in with the username and password specified as `MONGO_EXPRESS_LOGIN_USERNAME`and `MONGO_EXPRESS_LOGIN_PASSWORD`in the [.env](https://github.com/sergicanet9/go-hexagonal-api/blob/main/.env) file.

## ☁️ Live Environment
The API is deployed on a Google Kubernetes Engine (GKE) cluster, using Mongo Atlas as database.
<br />
The application is instrumented with the New Relic Go agent for APM and log forwarding.
<br/>
For public access, the Kubernetes Service is configured as ClusterIP, and traffic is routed to a custom domain using a Cloudflare tunnel, avoiding the need for a public IP on the cluster.
<br />
Check the Swagger at:
https://mongo-prod-go-hexagonal-api.sergicanet.com/swagger/index.html
<br/>
<br/>
NOTES:
<br/>
To minimize costs when the API is not in use, the cluster is configured to scale to zero nodes. This means that:
- The first access to the API after a period of inactivity may take a few minutes. This delay is due to GKE needing to provision and start a new virtual machine node to host the API.
- Once the node is active, the API will respond quickly to subsequent requests.


## ✍️ Author
Sergi Canet Vela

## ⚖️ License
This project is licensed under the terms of the MIT license.
