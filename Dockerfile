# Start from the latest golang base image
FROM golang:alpine

RUN GOCACHE=OFF

RUN go env -w GOPRIVATE=github.com/scanet9

RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

RUN git config --global url."https://golang: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@github.com:scanet9".insteadOf "https://github.com/scanet9"

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

#ENTRYPOINT ["/app"]

# Command to run the executable
CMD ["./main"]