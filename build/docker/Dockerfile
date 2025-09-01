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
COPY --from=builder /opt/go-hexagonal-api/proto/gen/openapiv2/*.json /opt/go-hexagonal-api/proto/gen/openapiv2/

WORKDIR /opt/go-hexagonal-api

ARG version
ENV v=$version

ARG environment
ENV env=$environment

ARG httpPort
ENV hp=$httpPort

ARG grpcPort
ENV gp=$grpcPort

ARG database
ENV db=$database

ARG dsn
ENV dsn=$dsn

ARG newrelicKey
ENV nrkey=$newrelicKey

CMD ["sh", "-c", "bin/main --ver $v --env $env --hport $hp --gport $gp --db $db --dsn $dsn --nrkey \"$nrkey\""]
