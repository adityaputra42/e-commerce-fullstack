# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy entire source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/main.go

# Runtime Stage
FROM alpine:3.23

WORKDIR /app

# Install runtime dependencies (including wget for healthcheck)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    postgresql-client \
    wget \
    curl \
    && rm -rf /var/cache/apk/*

# Set timezone
ENV TZ=Asia/Jakarta

# Copy binary from builder
COPY --from=builder /app/main .

# Copy entrypoint script
COPY entrypoint.sh .

# Copy .env file
COPY .env .

# Make scripts executable
RUN chmod +x /app/entrypoint.sh /app/main

# Expose port
EXPOSE 8080

# Run entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]