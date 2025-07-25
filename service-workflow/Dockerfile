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
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer ./cmd/consumer

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/service /app/service
COPY --from=builder /app/worker /app/worker
COPY --from=builder /app/consumer /app/consumer
COPY --from=builder /app/ticket_workflow.json /app/ticket_workflow.json

# Run as non-root user
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

# The application will be started with the appropriate command in the deployment 