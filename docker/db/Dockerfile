# Start from the latest golang base image
FROM golang:1.21-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=builder /app/main .

# Note: The .env file is not copied into the image for security reasons.
# Environment variables should be passed directly to the container.

CMD ["./main"]
