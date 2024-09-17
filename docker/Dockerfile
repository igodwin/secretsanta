FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o secretsanta ./cmd/cli

FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/secretsanta /app/

RUN chown -R appuser:appgroup /app

USER appuser

ENTRYPOINT ["/app/secretsanta"]
