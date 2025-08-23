# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=${VERSION} -s -w" \
    -a -installsuffix cgo \
    -o testdoc ./cmd/testdoc

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN addgroup -g 1001 -S testdoc && \
    adduser -u 1001 -S testdoc -G testdoc

# Set working directory
WORKDIR /workspace

# Copy binary from builder stage
COPY --from=builder /app/testdoc /usr/local/bin/testdoc

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Keep root user for volume access
# USER testdoc

# Set default entrypoint
ENTRYPOINT ["testdoc"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="TestDoc"
LABEL org.opencontainers.image.description="Automatic Go test documentation generator"
LABEL org.opencontainers.image.vendor="TestDoc Organization"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/seblex/testdoc"
