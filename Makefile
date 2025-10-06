.PHONY: all build test test-unit test-verbose test-short test-k8sclient clean help stubgen example

# Default target - runs ALL tests (unit + integration)
all: test build

# Help target
help:
	@echo "Available targets:"
	@echo "  all            - Run ALL tests (unit + k8sclient integration) and build everything (default)"
	@echo "  test           - Run ALL tests: unit tests + k8sclient integration test"
	@echo "  test-unit      - Run only unit tests with race detection and coverage"
	@echo "  test-verbose   - Run unit tests with verbose output (shows every test)"
	@echo "  test-short     - Run unit tests without race detection (faster)"
	@echo "  test-k8sclient - Run k8sclient example with Kind cluster (requires kind & kubectl)"
	@echo "  build          - Build all binaries"
	@echo "  stubgen        - Build stubgen code generator"
	@echo "  example        - Build example application"
	@echo "  clean          - Remove built binaries"
	@echo "  help           - Show this help message"

# Run ALL tests (unit + integration)
test: test-unit test-k8sclient
	@echo ""
	@echo "✓ ALL tests passed (unit + integration)"

# Run unit tests with race detection and coverage
test-unit:
	@echo "=== Running unit tests with race detection and coverage ==="
	go test -race -cover ./pkg/...
	@echo ""
	@echo "✓ Unit tests passed"

# Run all tests with verbose output (shows every test function)
test-verbose:
	@echo "=== Running all tests (VERBOSE MODE) ==="
	@echo ""
	@echo "📦 Package: pkg/glua"
	@go test -v -race -cover ./pkg/glua
	@echo ""
	@echo "📦 Package: pkg/modules/kubernetes"
	@go test -v -race -cover ./pkg/modules/kubernetes
	@echo ""
	@echo "📦 Package: pkg/stubgen"
	@go test -v -race -cover ./pkg/stubgen
	@echo ""
	@echo "✓ All tests passed"

# Run all tests without race detection (faster)
test-short:
	@echo "=== Running all tests (SHORT MODE - no race detection) ==="
	go test -cover ./pkg/...
	@echo ""
	@echo "✓ All tests passed"

# Run k8sclient example with Kind cluster
test-k8sclient:
	@echo "=== Running k8sclient example test with Kind cluster ==="
	@if ! command -v kind >/dev/null 2>&1; then \
		echo "Error: kind is not installed. Install from https://kind.sigs.k8s.io/"; \
		exit 1; \
	fi
	@if ! command -v kubectl >/dev/null 2>&1; then \
		echo "Error: kubectl is not installed. Install from https://kubernetes.io/docs/tasks/tools/"; \
		exit 1; \
	fi
	@cd example/k8sclient && ./run-test.sh

# Build all binaries
build: stubgen example
	@echo ""
	@echo "✓ All binaries built successfully"

# Build stubgen code generator
stubgen:
	@echo "=== Building stubgen ==="
	go build -o bin/stubgen ./cmd/stubgen
	@echo "✓ stubgen built -> bin/stubgen"

# Build example application
example:
	@echo "=== Building example ==="
	cd example && go build -o ../bin/example .
	@echo "✓ example built -> bin/example"

# Clean built binaries
clean:
	@echo "=== Cleaning built binaries ==="
	rm -rf bin/
	@echo "✓ Clean complete"
