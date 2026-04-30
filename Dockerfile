# ─── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /app/bin && CGO_ENABLED=0 GOOS=linux go build -o /app/bin/sluggo ./cmd/api

# ─── Runtime stage ────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/sluggo ./bin/sluggo

EXPOSE 8080
CMD ["./bin/sluggo"]
