.PHONY: test docs

up:
	export TAG=$(shell git branch --show-current | xargs basename) && docker-compose -f docker/docker-compose.yml -f docker/docker-compose.local.yml up -d --build
down:
	docker-compose -f docker/docker-compose.yml down
test:
	go test ./test -coverpkg=./... -coverprofile=test/coverage.out
cover:
	go tool cover -html=test/coverage.out
docs:
	go install github.com/swaggo/swag/cmd/swag@v1.7.0
	swag init -g cmd/main.go
goose-create:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir adapters/repositories/postgres/migrations/ create $${name:-<name>} sql
goose-up:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Connection string (format: host=XX port=XX dbname=XX user=XX password=XX sslmode=XX): " connection_string; \
	goose -dir adapters/repositories/postgres/migrations/ -table "public.goose_db_version" postgres "$${connection_string:-<connection_string>} " up
