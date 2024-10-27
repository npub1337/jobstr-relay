# Variables
BINARY_NAME=jobstr-relay
DB_NAME=jobstr.db
TEST_BINARY=test-relay
DB_DIR=db

# Default commands
.PHONY: all build run clean setup db-dir test test-ci run-dev stop-relay start-relay

all: setup

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME)

# Create the database directory with correct permissions
db-dir:
	@echo "Ensuring database directory exists with correct permissions..."
	@mkdir -p $(DB_DIR)
	@chmod 755 $(DB_DIR)

# Stop the relay if it's running
stop-relay:
	@echo "Checking if relay is running..."
	@if pgrep -x "$(BINARY_NAME)" > /dev/null; then \
		echo "Stopping existing relay..."; \
		pkill -x "$(BINARY_NAME)"; \
		sleep 2; \
	else \
		echo "No existing relay found."; \
	fi

# Start the relay
start-relay: stop-relay
	@echo "Starting relay..."
	@./$(BINARY_NAME) &
	@sleep 2
	@echo "Relay started."

# Run the relay
run: db-dir build start-relay

# Clean up generated files and database directory
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(TEST_BINARY)
	@rm -rf $(DB_DIR)
	@echo "Cleaned up successfully!"

# Set up and run the project
setup: clean db-dir build run

# Run tests locally with verbose output
test: clean build db-dir start-relay
	@echo "Running tests with verbose output..."
	@go test ./test -v
	@echo "Stopping relay..."
	@$(MAKE) stop-relay
	@echo "Cleaning up after test..."
	@$(MAKE) clean

# Run tests for CI with verbose output
test-ci: clean build db-dir start-relay
	@echo "Running tests with verbose output..."
	@go test ./test -v
	@echo "Stopping relay..."
	@$(MAKE) stop-relay
	@echo "Cleaning up after test..."
	@$(MAKE) clean

# Run locally for development
run-dev: db-dir build start-relay
	@echo "Running $(BINARY_NAME) in development mode..."
	@tail -f /dev/null
