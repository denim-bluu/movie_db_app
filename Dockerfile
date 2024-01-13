FROM golang:1.21-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the wait-for-it script
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Command to run the binary
CMD ["./main"]
