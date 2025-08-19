# Build stage
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./cmd/server

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-service .
COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./order-service"]