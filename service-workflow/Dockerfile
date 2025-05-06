# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/service
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/service /app/service
COPY --from=builder /app/worker /app/worker

# Run as non-root user
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

# The application will be started with the appropriate command in the deployment 