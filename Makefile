include .env

.PHONY: test

up:
	rm -f mongo.keyfile
	openssl rand -base64 24 > mongo.keyfile
	chmod 400 mongo.keyfile
	chown :admin mongo.keyfile
	docker-compose up -d --build
	@echo "Mongo Swagger:    http://localhost:${HOST_PORT_MONGOAPI}/swagger/index.html"
	@echo "Postgres Swagger: http://localhost:${HOST_PORT_POSTGRESAPI}/swagger/index.html"
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
