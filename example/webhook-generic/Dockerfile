# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /workspace

# Copy go mod files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy source code
COPY ../../pkg ./pkg
COPY . ./example/webhook-generic

# Build
WORKDIR /workspace/example/webhook-generic
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o webhook .

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy binary
COPY --from=builder /workspace/example/webhook-generic/webhook /webhook

# Copy Lua scripts
COPY --from=builder /workspace/example/webhook-generic/scripts /etc/webhook/scripts

USER 65532:65532

ENTRYPOINT ["/webhook"]
CMD ["-address=:8443", "-cert=/etc/webhook/certs/tls.crt", "-key=/etc/webhook/certs/tls.key", "-scripts=/etc/webhook/scripts"]
