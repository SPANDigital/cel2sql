# Makefile for cel2sql project

.PHONY: build test lint fmt clean help install-tools deps vuln-check

# Build the project
build:
	go build -v ./...

# Run tests
test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -cache -testcache -modcache

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

# Download dependencies
deps:
	go mod download
	go mod verify

# Run vulnerability check
vuln-check:
	govulncheck ./...

# Run all checks (used in CI)
ci: fmt lint test vuln-check

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Run the example
example:
	go run examples/postgresql_example.go

# Generate documentation
docs:
	go doc -all ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the project"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linting"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  install-tools - Install development tools"
	@echo "  deps          - Download and verify dependencies"
	@echo "  vuln-check    - Run vulnerability check"
	@echo "  ci            - Run all checks (used in CI)"
	@echo "  update-deps   - Update dependencies"
	@echo "  example       - Run the PostgreSQL example"
	@echo "  docs          - Generate documentation"
	@echo "  help          - Show this help message"
