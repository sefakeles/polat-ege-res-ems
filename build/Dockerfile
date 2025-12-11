FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-s -w -X main.Version=docker -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -a -installsuffix cgo -o ems ./cmd/ems

# Final stage
FROM alpine:latest

# Install runtime dependencies including curl for health check
RUN apk --no-cache add ca-certificates sqlite curl tzdata

# Create app user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/ems .

# Copy configuration files
COPY configs/ ./configs/

# Create directories with proper permissions
RUN mkdir -p data logs && \
    chown -R appuser:appuser /app

# Switch to app user
USER appuser

# Expose port
EXPOSE 8080

# Health check using curl instead of wget
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/status || exit 1

# Run the application
CMD ["./ems"]
