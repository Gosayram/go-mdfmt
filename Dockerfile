# Build stage v1.24.4-alpine3.22
FROM golang:1.24.4-alpine3.22@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/Gosayram/go-mdfmt/internal/version.Version=docker" \
    -o mdfmt ./cmd/mdfmt

# Final stage
FROM scratch

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/mdfmt /usr/local/bin/mdfmt

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/mdfmt"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="mdfmt"
LABEL org.opencontainers.image.description="Fast, reliable Markdown formatter"
LABEL org.opencontainers.image.source="https://github.com/Gosayram/go-mdfmt"
LABEL org.opencontainers.image.licenses="MIT" 