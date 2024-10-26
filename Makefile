# Variables
BINARY_NAME=jobstr-relay
DB_NAME=jobstr.db
TEST_BINARY=test-relay

# Default commands
.PHONY: all build run clean setup db-dir test test-ci

all: setup

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME)

# Create the database directory with correct permissions
db-dir:
	@echo "Ensuring database directory exists with correct permissions..."
	@mkdir -p db
	@chmod 755 db

# Run the relay
run: db-dir build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Clean up generated files
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(TEST_BINARY)
	@rm -f db/$(DB_NAME)
	@echo "Cleaned up successfully!"

# Set up and run the project
setup: clean db-dir build run

# Run tests locally
test: build
	@echo "Running relay in background..."
	@./$(BINARY_NAME) &
	@echo "Waiting for relay to start..."
	@sleep 2
	@echo "Running tests..."
	@go test ./test -v
	@echo "Stopping relay..."
	@pkill -f $(BINARY_NAME)
	@echo "Cleaning up after test..."
	@$(MAKE) clean

# Run tests for CI
test-ci: build
	@echo "Running relay in background..."
	@./$(BINARY_NAME) &
	@echo "Waiting for relay to start..."
	@sleep 2
	@echo "Running tests..."
	@go test ./test -v
	@echo "Stopping relay..."
	@pkill -f $(BINARY_NAME)

# Run locally for development
run-dev: db-dir build
	@echo "Running $(BINARY_NAME) in development mode..."
	@./$(BINARY_NAME)
