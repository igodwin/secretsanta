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

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application for target architecture
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} GOARM=${TARGETVARIANT#v} \
    go build -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o secretsanta \
    ./cmd/cli

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest AS production

# Build arguments for metadata
ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S secretsanta && \
    adduser -u 1001 -S secretsanta -G secretsanta

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/secretsanta .

# Copy config template
COPY --from=builder /app/configs/secretsanta.config.template ./config.template

# Change ownership to non-root user
RUN chown -R secretsanta:secretsanta /app

# Switch to non-root user
USER secretsanta

# Default command
ENTRYPOINT ["./secretsanta"]
CMD ["participants.json"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep secretsanta || exit 1

# Metadata with build info
LABEL maintainer="igodwin" \
      description="Secret Santa application with gRPC notifier integration" \
      version="${VERSION:-1.0.0}" \
      org.opencontainers.image.title="Secret Santa" \
      org.opencontainers.image.description="Secret Santa gift exchange application with gRPC notifier integration" \
      org.opencontainers.image.authors="igodwin" \
      org.opencontainers.image.vendor="igodwin" \
      org.opencontainers.image.version="${VERSION:-1.0.0}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.source="https://github.com/igodwin/secretsanta" \
      org.opencontainers.image.documentation="https://github.com/igodwin/secretsanta" \
      org.opencontainers.image.licenses="MIT" \
      org.label-schema.build-date="${BUILD_DATE}" \
      org.label-schema.vcs-ref="${VCS_REF}" \
      org.label-schema.version="${VERSION:-1.0.0}" \
      org.label-schema.schema-version="1.0" \
      architecture="${TARGETARCH}" \
      os="${TARGETOS}"