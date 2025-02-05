# Application name
APP_NAME = smtp-proxy-api
# Go binary name
BINARY_NAME = $(APP_NAME)
# Go module name
MODULE_NAME = github.com/sinland/$(APP_NAME)
# Go version
GO_VERSION = 1.23
# Build output directory
BUILD_DIR = bin

# Default target
all: build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Build the application
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(APP_NAME)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Lint the code
lint:
	@echo "Linting..."
	@golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Check Go version
check-go-version:
	@echo "Checking Go version..."
	@go version | grep -q "go$(GO_VERSION)" || (echo "Go version $(GO_VERSION) is required" && exit 1)


# Deploy
deploy:
	@scp bin/rexpatgame-backend admin@62.84.117.178:/opt/apps/rexpatgame-backend

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  deps         - Install dependencies"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint         - Lint the code"
	@echo "  fmt          - Format the code"
	@echo "  clean        - Clean build artifacts"
	@echo "  check-go-version - Check if the correct Go version is installed"
	@echo "  help         - Show this help message"

.PHONY: all deps build run test test-coverage lint fmt clean check-go-version help