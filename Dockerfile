# Start from the latest golang base image
FROM golang:alpine

RUN GOCACHE=OFF

RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /app/go-mongo-restapi

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o ./out/go-mongo-restapi cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./out/go-mongo-restapi"]