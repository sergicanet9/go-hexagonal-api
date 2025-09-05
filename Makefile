include .env

.PHONY: test

up:
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	docker-compose up -d --build
	@echo ""
	@echo "👉 Mongo API Swagger:    http://localhost:${HOST_HTTP_PORT_MONGOAPI}/swagger/index.html"
	@echo "👉 Mongo API gRPC UI: 	 http://localhost:${HOST_HTTP_PORT_MONGOAPI}/grpcui/"
	@echo "👉 Mongo Express:        http://localhost:${MONGO_EXPRESS_HOST_PORT}"
	@echo ""
	@echo "👉 Postgres API Swagger: http://localhost:${HOST_HTTP_PORT_POSTGRESAPI}/swagger/index.html"
	@echo "👉 Postgres API gRPC UI: http://localhost:${HOST_HTTP_PORT_POSTGRESAPI}/grpcui/"
	@echo "👉 PgAdmin:              http://localhost:${PGADMIN_HOST_PORT}"
	@echo ""
down:
	docker-compose down
test-unit:
	go test -race $(shell go list ./... | grep -v /test) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
test-integration:
	go test -race test/integration/*.go
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --dir=core/ports --all --output=test/mocks
goose:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir infrastructure/postgres/migrations/ create $$name sql
protos:
	cd proto && buf dep update
	cd proto && buf generate --clean