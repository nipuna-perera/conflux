# Multi-stage Docker build for Go backend
# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Stage 2: Create minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
