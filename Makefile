.PHONY: test

mongo-up:
	export VERSION=$(shell git branch --show-current | xargs basename) && export ENV=local && export PORT=8080 && export DB=mongo && docker-compose up -d --build
postgres-up:
	export VERSION=$(shell git branch --show-current | xargs basename) && export ENV=local && export PORT=8080 && export DB=postgres && docker-compose up -d --build
down:
	export VERSION= && export ENV= && export PORT=0 && export DB= && docker-compose down
test:
	go test ./test -coverpkg=./... -coverprofile=test/coverage.out
cover:
	go tool cover -html=test/coverage.out
docs:
	go install github.com/swaggo/swag/cmd/swag@v1.7.0
	swag init -g cmd/main.go -o app/docs
goose-create:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir adapters/postgres/migrations/ create $${name:-<name>} sql
goose-up:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Connection string (format: host=XX port=XX dbname=XX user=XX password=XX sslmode=XX): " connection_string; \
	goose -dir adapters/postgres/migrations/ -table "public.goose_db_version" postgres "$${connection_string:-<connection_string>} " up
