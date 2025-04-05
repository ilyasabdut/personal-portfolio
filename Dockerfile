# -------- Build Stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Add build tools
RUN apk add --no-cache git gcc musl-dev

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN rm -rf ./sqlite/data.db

RUN go build -o server .

# -------- Runtime Stage --------
FROM alpine:latest

WORKDIR /app

# Copy the built binary only
COPY --from=builder /app/server .

# Copy static files and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

EXPOSE 8080
CMD ["./server"]
