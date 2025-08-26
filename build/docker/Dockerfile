FROM golang:alpine AS builder

WORKDIR /opt/go-hexagonal-api

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o bin/main cmd/main.go

FROM alpine:latest

COPY --from=builder /opt/go-hexagonal-api/bin/main /opt/go-hexagonal-api/bin/main
COPY --from=builder /opt/go-hexagonal-api/config/*.json /opt/go-hexagonal-api/config/
COPY --from=builder /opt/go-hexagonal-api/infrastructure/postgres/migrations/*.sql /opt/go-hexagonal-api/infrastructure/postgres/migrations/

WORKDIR /opt/go-hexagonal-api

ARG version
ENV v $version

ENV env $environment
ENV p $port
ENV db $database
ENV dsn $dsn

CMD ["sh", "-c", "bin/main --ver $v --env $env --port $p --db $db --dsn $dsn"]