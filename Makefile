include .env

.PHONY: test

up:
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	docker-compose up -d --build
	@echo ""
	@echo "🍃 Mongo API"
	@echo "    👉 Swagger UI:    http://localhost:${HOST_HTTP_PORT_MONGOAPI}/v1/swagger/index.html"
	@echo "    👉 gRPC UI:       http://localhost:${HOST_HTTP_PORT_MONGOAPI}/v1/grpcui/"
	@echo "    👉 Mongo Express: http://localhost:${MONGO_EXPRESS_HOST_PORT}"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost:${HOST_HTTP_PORT_MONGOAPI}/v1/health"
	@echo "        grpcurl -plaintext localhost:${HOST_GRPC_PORT_MONGOAPI} health.HealthService/HealthCheck"
	@echo ""
	@echo "🐘 Postgres API"
	@echo "    👉 Swagger UI:    http://localhost:${HOST_HTTP_PORT_POSTGRESAPI}/v1/swagger/index.html"
	@echo "    👉 gRPC UI:       http://localhost:${HOST_HTTP_PORT_POSTGRESAPI}/v1/grpcui/"
	@echo "    👉 pgAdmin:       http://localhost:${PGADMIN_HOST_PORT}"
	@echo "    🔧 Command examples:"
	@echo "        curl http://localhost:${HOST_HTTP_PORT_POSTGRESAPI}/v1/health"
	@echo "        grpcurl -plaintext localhost:${HOST_GRPC_PORT_POSTGRESAPI} health.HealthService/HealthCheck"
	@echo ""
down:
	docker-compose down
test-unit:
	go test -race $(shell go list ./... | grep -v /test | grep -v /pb) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
test-integration:
	go test -race test/integration/*.go
protos:
	cd proto && buf dep update
	cd proto && buf generate --clean
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --dir=core/ports --all --output=test/mocks
goose:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir infrastructure/postgres/migrations/ create $$name sql
