# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the webhook binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o webhook ./example/webhook

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/webhook .

# Copy Lua mutation script
COPY example/webhook/scripts/mutate.lua /etc/webhook/scripts/mutate.lua

# Create directory for certs (will be mounted as volume)
RUN mkdir -p /etc/webhook/certs

EXPOSE 8443

ENTRYPOINT ["/app/webhook"]
CMD ["--address=:8443", "--cert=/etc/webhook/certs/tls.crt", "--key=/etc/webhook/certs/tls.key", "--script=/etc/webhook/scripts/mutate.lua"]
