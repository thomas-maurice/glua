.PHONY: all build test test-verbose test-short clean help stubgen example

# Default target
all: test build

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Run tests and build everything (default)"
	@echo "  test         - Run all tests with race detection and coverage"
	@echo "  test-verbose - Run all tests with verbose output (shows every test)"
	@echo "  test-short   - Run all tests without race detection (faster)"
	@echo "  build        - Build all binaries"
	@echo "  stubgen      - Build stubgen code generator"
	@echo "  example      - Build example application"
	@echo "  clean        - Remove built binaries"
	@echo "  help         - Show this help message"

# Run all tests with race detection and coverage
test:
	@echo "=== Running all tests with race detection and coverage ==="
	go test -v -race -cover ./pkg/...
	@echo ""
	@echo "✓ All tests passed"

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
