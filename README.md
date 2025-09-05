# go-hexagonal-api
![CI](https://github.com/sergicanet9/go-hexagonal-api/actions/workflows/ci.yml/badge.svg)
![CD](https://github.com/sergicanet9/go-hexagonal-api/actions/workflows/cd.yml/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-73.2%25-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergicanet9/go-hexagonal-api.svg)](https://pkg.go.dev/github.com/sergicanet9/go-hexagonal-api)

A robust gRPC + REST API for user management built with **Go** and implementing the **Hexagonal Architecture** (Ports & Adapters) pattern. It makes use of [scv-go-tools](https://github.com/sergicanet9/scv-go-tools) library.
<br />
The gRPC handlers are automatically exposed as REST endpoints, so the same functionality is available over HTTP without duplicating code.
<br />
The API is designed to work seamlessly with either a MongoDB or PostgreSQL database instance, using the same business logic and handlers.

## 🚀 Features
- **Hexagonal Architecture**: Clear separation of concerns with transport, business logic and repository layers.
- **gRPC + gRPC-Gateway**: gRPC API implementation, with automatically generated REST endpoints from the gRPC handlers via gRPC-Gateway.
- **Database Agnostic**: Decoupled repository adapters allow injecting a MongoDB or PostgreSQL storage without changing core logic.
- **Authentication & Authorization**: Implements JWT authentication and claim-based authorization for secure endpoints.
- **Asynchronous Processes**: Go routines management with built in processes for periodically health checking connectivity with HTTP and gRPC servers.
- **Testing**: Comprehensive unit tests with code coverage and integration tests for the happy path.
- **Developer Experience**: Built-in Makefile, Swagger UI, gRPC UI, pgAdmin, and mongo-express.
- **Lifecycle Management**: Multi-environment support with config files, Dockerfile and docker-compose, CI/CD pipelines with GitHub Actions, Kubernetes deployment and New Relic observability.

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

Check the console output for Swagger UI, gRPC UI, database UIs, and HTTP and gRPC command examples.

To stop and remove all containers, run:
```
make down
```

### Run it with command line
Run a single API instance with command-line arguments with the following command:
```
    go run cmd/main.go --ver={version} --env={environment} --hport={http_port} --gport={grpc_port} --db={database} --dsn={dsn} --nrkey={newrelic_key}
```
or:
```
go build cmd/main.go
 ./main --ver={version} --env={environment} --hport={http_port} --gport={grpc_port} --db={database} --dsn={dsn} --nrkey={newrelic_key}
```
Provide the desired values for: `{version}`, `{environment}`, `{http_port}`, `{grpc_port}`, `{database}`, `{dsn}`.
<br />
The `--nrkey` flag and its value `{newrelic_key}` are optional and can be omitted if you do not want to configure New Relic observability.
<br />
Then open:
* Swagger UI: `http://localhost:{http_port}/swagger/index.html`
* gRPC UI: `http://localhost:{http_port}/grpcui/`
<br />

NOTES:
- The target database container needs to be up and running (run `make up`).

### Debug it with VS Code
The project includes debugging profiles in [launch.json](https://github.com/sergicanet9/go-hexagonal-api/blob/main/.vscode/launch.json) for both MongoDB and PostgreSQL setups. Simply select the desired configuration in the VS Code debugger and run it.
<br />
Then open:
* Swagger UI: `http://localhost:{http_port}/swagger/index.html`
* gRPC UI: `http://localhost:{http_port}/grpcui/`
<br />

NOTES:
- The target database container needs to be up and running (run `make up`).

## 📦 API Endpoints
### Public Routes
These endpoints do not require authentication.

| HTTP Endpoint         | gRPC Method                         | Description                                   |
| :-------------------- | :---------------------------------- | :-------------------------------------------- |
| GET `/health`         | `health.HealthService.HealthCheck`  | Performs a health check.                      |
| POST `/v1/users`      | `user.UserService.Create`           | Creates a new user.                           |
| POST `/v1/users/many` | `user.UserService.CreateMany`       | Creates multiple users.                       |
| POST `/v1/users/login`| `user.UserService.Login`            | Authenticates a user and returns a JWT token. |

### Protected Routes
These endpoints require a valid JWT in the Authorization header, formatted as `Bearer {token}`.
* For HTTP, include it as `Authorization` header.
* For gRPC, include it in the metadata with the key `authorization`.

| HTTP Endpoint                  | gRPC Method                    | Description                   |
| :----------------------------- | :----------------------------- | :---------------------------- |
| GET `/v1/users`                | `user.UserService.GetAll`      | Retrieves all users.          |
| GET `/v1/users/email/{email}`  | `user.UserService.GetByEmail`  | Retrieves a user by email.    |
| GET `/v1/users/{id}`           | `user.UserService.GetByID`     | Retrieves a user by ID.       |
| PATCH `/v1/users/{id}`         | `user.UserService.Update`      | Updates a user's information. |
| GET `/v1/claims`               | `user.UserService.GetClaims`   | Returns all claims.           |

### Admin Routes
These endpoints require a valid JWT, formatted as `Bearer {token}` and containing the `admin` claim.
* For HTTP, include it as `Authorization` header.
* For gRPC, include it in the metadata with the key `authorization`.

| HTTP Endpoint           | gRPC Method                | Description           |
| :---------------------- | :------------------------- | :-------------------- |
| DELETE `/v1/users/{id}` | `user.UserService.Delete`  | Deletes a user by ID. |

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
<br />

 NOTES:
- Docker is required for running integration tests.

## 🛠️ Developer Commands 
### (Re)Generate gRPC stubs and Swagger documentation
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
The API is deployed on a Google Kubernetes Engine (GKE) cluster, using Mongo Atlas as database, New Relic Go agent for APM and log forwarding, and a Cloudflare tunnel for public access through HTTP.<br/>
* Swagger UI: https://mongo-prod-go-hexagonal-api.sergicanet.com/swagger/index.html
* gRPC UI: https://mongo-prod-go-hexagonal-api.sergicanet.com/grpcui/
<br />

NOTES:
- The gRPC UI is publically exposed through HTTP. However, direct gRPC calls (e.g., using `grpcurl`) are not allowed outside the cluster.

## ✍️ Author
Sergi Canet Vela

## ⚖️ License
This project is licensed under the terms of the MIT license.
