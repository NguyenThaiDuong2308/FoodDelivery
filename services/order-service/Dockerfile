# 1. Build stage
FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service cmd/main.go

# 2. Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/order-service .
COPY config/ ./config/
EXPOSE 8083
CMD ["./order-service"]

