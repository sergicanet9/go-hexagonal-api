include .env

.PHONY: test

up:
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	docker-compose up -d --build
	@echo "Mongo API Swagger:    http://localhost:${HOST_PORT_MONGOAPI}/swagger/index.html"
	@echo "Postgres API Swagger: http://localhost:${HOST_PORT_POSTGRESAPI}/swagger/index.html"
	@echo "Mongo Express:        http://localhost:${MONGO_EXPRESS_HOST_PORT}"
	@echo "PgAdmin:              http://localhost:${PGADMIN_HOST_PORT}"
down:
	docker-compose down
test-unit:
	go test -race $(shell go list ./... | grep -v /test) -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
test-integration:
	go test -race test/integration/*.go
swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.7.0
	swag init -g cmd/main.go -o app/docs
mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --dir=core/ports --all --output=test/mocks
goose:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir infrastructure/postgres/migrations/ create $$name sql
protos:
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o proto/google/api/annotations.proto
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o proto/google/api/http.proto
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
# 	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
# 	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16
	cd proto && protoc --go_out=protogen --go_opt=paths=source_relative \
	./*.proto
