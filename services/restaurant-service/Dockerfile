# 1. Build stage
FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o restaurant-service cmd/main.go

# 2. Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/restaurant-service .
COPY config/ ./config/
EXPOSE 8081
CMD ["./restaurant-service"]