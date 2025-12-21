# ================================
# Build Stage
# ================================
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates 
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/api \
    ./cmd/main.go

# ================================
# Production Stage
# ================================
FROM scratch

# Import CA certificates for HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /app/api /api

# Copy migrations 
COPY --from=builder /app/cmd/migrations /migrations

# Expose the application port
EXPOSE 8888

# Set environment defaults 
ENV PORT=":8888" \
    ENVIRONMENT="production" \
    LOG_FORMAT="json" \
    LOG_LEVEL="info"

# Run the binary
ENTRYPOINT ["/api"]
