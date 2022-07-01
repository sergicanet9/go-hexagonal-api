# Start from the latest golang base image
FROM golang:alpine

RUN GOCACHE=OFF

RUN apk add git

ARG version
ENV v $version
ARG environment
ENV env $environment
ARG port
ENV p $port
ARG database
ENV db $database

# Set the Current Working Directory inside the container
WORKDIR /app/go-hexagonal-api

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o ./opt/go-hexagonal-api cmd/main.go

# Command to run the executable
CMD ["sh", "-c", "./opt/go-hexagonal-api -v $v -env $env -p $p -db $db"]