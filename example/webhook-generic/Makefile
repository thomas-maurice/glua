.PHONY: build run clean docker-build docker-push

BINARY=webhook-generic
IMAGE=ghcr.io/thomas-maurice/glua-webhook-generic
TAG?=latest

build:
	@echo "Building $(BINARY)..."
	go build -o bin/$(BINARY) .

run: build
	@echo "Running $(BINARY)..."
	./bin/$(BINARY) \
		-address=:8443 \
		-cert="" \
		-key="" \
		-scripts=./scripts \
		-enable-nodes=true \
		-enable-pods=true

clean:
	@echo "Cleaning..."
	rm -rf bin/

docker-build:
	@echo "Building Docker image $(IMAGE):$(TAG)..."
	docker build -t $(IMAGE):$(TAG) .

docker-push: docker-build
	@echo "Pushing Docker image $(IMAGE):$(TAG)..."
	docker push $(IMAGE):$(TAG)

help:
	@echo "Available targets:"
	@echo "  build        - Build the webhook binary"
	@echo "  run          - Build and run locally (without TLS)"
	@echo "  clean        - Remove built binaries"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Build and push Docker image"
	@echo ""
	@echo "Variables:"
	@echo "  IMAGE        - Docker image name (default: $(IMAGE))"
	@echo "  TAG          - Docker image tag (default: $(TAG))"
