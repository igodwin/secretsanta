# Multi-stage, multi-architecture Dockerfile for Secret Santa application
# syntax=docker/dockerfile:1

# Development stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS dev
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["tail", "-f", "/dev/null"]

# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Build arguments
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the web server for target architecture with version info
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} GOARM=${TARGETVARIANT#v} \
    go build -ldflags="-w -s -extldflags '-static' -X 'main.Version=${VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'" \
    -a -installsuffix cgo \
    -o secretsanta-web \
    ./cmd/web

# Final stage
FROM alpine:latest AS production

# Build arguments for metadata
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S secretsanta && \
    adduser -u 1001 -S secretsanta -G secretsanta

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/secretsanta-web .

# Copy static web assets
COPY --from=builder /app/internal/web/static ./internal/web/static

# Copy config template
COPY --from=builder /app/configs/config.yaml ./config.template

# Change ownership to non-root user
RUN chown -R secretsanta:secretsanta /app

# Switch to non-root user
USER secretsanta

# Expose web server port
EXPOSE 8080

# Default command
ENTRYPOINT ["./secretsanta-web"]
CMD ["-addr", ":8080"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1

# Metadata with build info
LABEL maintainer="igodwin" \
      description="Secret Santa web application with gRPC notifier integration" \
      version="${VERSION:-1.0.0}" \
      org.opencontainers.image.title="Secret Santa Web" \
      org.opencontainers.image.description="Secret Santa web server with REST API and UI" \
      org.opencontainers.image.authors="igodwin" \
      org.opencontainers.image.vendor="igodwin" \
      org.opencontainers.image.version="${VERSION:-1.0.0}" \
      org.opencontainers.image.created="${BUILD_TIME}}" \
      org.opencontainers.image.revision="${GIT_COMMIT}}" \
      org.opencontainers.image.source="https://github.com/igodwin/secretsanta" \
      org.opencontainers.image.documentation="https://github.com/igodwin/secretsanta" \
      org.opencontainers.image.licenses="MIT" \
      org.label-schema.build-time="${BUILD_TIME}}" \
      org.label-schema.git-commit="${GIT_COMMIT}}" \
      org.label-schema.version="${VERSION:-1.0.0}" \
      org.label-schema.schema-version="1.0" \
      architecture="${TARGETARCH}" \
      os="${TARGETOS}"