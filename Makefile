.PHONY: test

mongo-up:
	export VERSION=$(shell git branch --show-current | xargs basename) && export ENV=local && export PORT=8080 && export DB=mongo && docker-compose up -d --build
postgres-up:
	export VERSION=$(shell git branch --show-current | xargs basename) && export ENV=local && export PORT=8080 && export DB=postgres && docker-compose up -d --build
down:
	export VERSION= && export ENV= && export PORT=0 && export DB= && docker-compose down
test:
	go test -race ./... -coverpkg=./... -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.7.0
	swag init -g cmd/main.go -o app/docs
goose:
	go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir db/postgres/migrations/ create $${name:-<name>} sql
