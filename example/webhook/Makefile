.PHONY: help build docker kind cert-manager install uninstall load-image test clean

CLUSTER_NAME ?= glua-webhook-test
IMAGE_NAME ?= glua-webhook
IMAGE_TAG ?= latest
NAMESPACE ?= glua-webhook
HELM_RELEASE ?= glua-webhook

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the webhook binary
	@echo "=== Building webhook binary ==="
	cd ../.. && go build -o example/webhook/bin/webhook ./example/webhook
	@echo "✓ Binary built -> bin/webhook"

docker: ## Build Docker image
	@echo "=== Building Docker image ==="
	cd ../.. && docker build -f example/webhook/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "✓ Docker image built: $(IMAGE_NAME):$(IMAGE_TAG)"

kind: ## Create Kind cluster
	@echo "=== Creating Kind cluster ==="
	@if kind get clusters | grep -q "^$(CLUSTER_NAME)$$"; then \
		echo "Cluster $(CLUSTER_NAME) already exists"; \
	else \
		kind create cluster --name $(CLUSTER_NAME); \
		echo "✓ Kind cluster created: $(CLUSTER_NAME)"; \
	fi

kind-delete: ## Delete Kind cluster
	@echo "=== Deleting Kind cluster ==="
	kind delete cluster --name $(CLUSTER_NAME)
	@echo "✓ Kind cluster deleted: $(CLUSTER_NAME)"

load-image: docker ## Load Docker image into Kind cluster
	@echo "=== Loading image into Kind cluster ==="
	kind load docker-image $(IMAGE_NAME):$(IMAGE_TAG) --name $(CLUSTER_NAME)
	@echo "✓ Image loaded into cluster"

cert-manager: ## Install cert-manager
	@echo "=== Installing cert-manager ==="
	@if kubectl get namespace cert-manager >/dev/null 2>&1; then \
		echo "cert-manager namespace already exists"; \
	else \
		helm repo add jetstack https://charts.jetstack.io --force-update; \
		helm repo update; \
		helm install cert-manager jetstack/cert-manager \
			--namespace cert-manager \
			--create-namespace \
			--set crds.enabled=true; \
		echo "Waiting for cert-manager to be ready..."; \
		kubectl wait --for=condition=available --timeout=300s \
			deployment/cert-manager -n cert-manager; \
		kubectl wait --for=condition=available --timeout=300s \
			deployment/cert-manager-webhook -n cert-manager; \
		echo "✓ cert-manager installed and ready"; \
	fi

install: load-image cert-manager ## Install the webhook using Helm
	@echo "=== Installing webhook via Helm ==="
	@kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install $(HELM_RELEASE) ./charts/glua-webhook \
		--namespace $(NAMESPACE) \
		--set image.repository=$(IMAGE_NAME) \
		--set image.tag=$(IMAGE_TAG) \
		--set image.pullPolicy=IfNotPresent \
		--wait
	@echo "✓ Webhook installed"

uninstall: ## Uninstall the webhook
	@echo "=== Uninstalling webhook ==="
	@helm uninstall $(HELM_RELEASE) --namespace $(NAMESPACE) || true
	@kubectl delete namespace $(NAMESPACE) --ignore-not-found
	@echo "✓ Webhook uninstalled"

test: ## Run test pod to verify mutation
	@bash test-webhook.sh

logs: ## Show webhook logs
	@kubectl logs -n $(NAMESPACE) -l app=glua-webhook --tail=100 -f

clean: ## Clean up local build artifacts
	@echo "=== Cleaning build artifacts ==="
	@rm -rf bin/
	@echo "✓ Cleaned"

all: kind install test ## Create cluster, install webhook, and run test
