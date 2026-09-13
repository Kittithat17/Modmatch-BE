# ===== Build stage =====
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download deps first so this layer is cached when only code changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# ===== Run stage =====
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
