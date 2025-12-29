# 1. Build Stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install build dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -o main: name the output file "main"
# ./cmd/api: points to where your main.go lives
RUN go build -o main ./cmd/api

# 2. Run Stage (Small image)
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env . 

# Expose port
EXPOSE 3000

# Run
CMD ["./main"]