# Name of the output binary
BINARY_NAME=task

# Directory for storing the compiled binary
BIN_DIR=bin

# Path to the main package
MAIN_PACKAGE=./main.go


# Default target: Build the binary
build:
	@echo "Building the CLI..."
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete! Executable is at ./$(BIN_DIR)/$(BINARY_NAME)"


# Clean the bin directory (delete the executable)
clean:
	@echo "Cleaning up..."
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "Cleanup complete."

# Format and lint the code
lint:
	@echo "Running code checks..."
	go fmt $(MAIN_PACKAGE)
	go vet $(MAIN_PACKAGE)
	@echo "Code checks complete."
