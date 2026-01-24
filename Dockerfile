# -------- Build Stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Add build tools
RUN apk add --no-cache git gcc musl-dev

# Copy and download dependencies
COPY . .

RUN go mod download

# Copy source and build
RUN rm -rf ./sqlite/data.db

RUN go build -o server ./cmd/server

# -------- Runtime Stage --------
FROM alpine:latest

WORKDIR /app

# Copy the built binary only
COPY --from=builder /app/server .

# Copy static files and templates
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates

EXPOSE 8081
CMD ["./server"]
