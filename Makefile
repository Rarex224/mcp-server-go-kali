.PHONY: build test test-integration clean run

IMAGE_NAME = backless-mcp-server:latest

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run Go unit tests locally
test:
	go test -v ./...

# Run the automated integration test script
test-integration: build
	./run_tests.sh

# Run the Docker container interactively
run: build
	docker run -i --rm $(IMAGE_NAME)

# Clean up Go build cache and tidy modules
clean:
	go clean
	go mod tidy
