# Start from the latest golang base image
FROM golang:1.21-bullseye as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Start a new stage from scratch
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /migrations

# Copy migration files
COPY db/migrations /migrations

# Copy the migrate tool from the builder stage
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Command to run the migrations
CMD migrate -path=/migrations -database ${DB_URL} up