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

