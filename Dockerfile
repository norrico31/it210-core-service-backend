# Use the official Go image for building the application
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git (if needed)
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application
COPY . .

# Clean up go.mod and go.sum
RUN go mod tidy

# Build the Go applications
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/seed ./cmd/seed/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/core-service ./cmd/main.go

# Use a lightweight image for the final stage
FROM alpine:latest
RUN apk add --no-cache netcat-openbsd
RUN apk add --no-cache libc6-compat

# Set the working directory for the final image
WORKDIR /app/core-service

# Copy the binaries from the builder stage
COPY --from=builder /app/core-service /app/core-service/
COPY --from=builder /app/migrate /app/core-service/migrate
COPY --from=builder /app/seed /app/core-service/seed
COPY entrypoint.sh /app/core-service/

# Copy the migration SQL files
COPY --from=builder /app/cmd/migrate/migrations /app/core-service/cmd/migrate/migrations

RUN chmod +x /app/core-service/entrypoint.sh /app/core-service/migrate /app/core-service/seed /app/core-service/core-service

EXPOSE 80 8080

ENTRYPOINT ["/app/core-service/entrypoint.sh"]
