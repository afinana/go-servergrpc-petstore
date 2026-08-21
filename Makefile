.PHONY: all build test test-unit test-coverage run docker-build docker-run compose-up compose-down clean tidy help

BINARY_NAME=petstore-server
DOCKER_IMAGE=go-servergrpc-petstore
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

all: build

## build: Compile the gRPC server binary
build:
	@echo "==> Building $(BINARY_NAME)..."
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) main.go

## run: Run the server locally
run:
	@echo "==> Running $(BINARY_NAME)..."
	go run main.go

## test: Run all tests (unit + integration if Mongo is up)
test:
	@echo "==> Running tests..."
	go test -v -race ./...

## test-coverage: Run tests and generate HTML coverage report
test-coverage:
	@echo "==> Generating test coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"

## tidy: Download and tidy Go dependencies
tidy:
	@echo "==> Tidying Go modules..."
	go mod tidy
	go mod verify

## vendor: Download dependencies into vendor folder for offline/container builds
vendor:
	@echo "==> Vendoring Go modules..."
	go mod vendor

## docker-build: Build the Docker image
docker-build: vendor
	@echo "==> Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE):latest .

## docker-run: Run the containerized server
docker-run:
	@echo "==> Running Docker container..."
	docker run --rm -p 8090:8090 --name $(BINARY_NAME) $(DOCKER_IMAGE):latest

## compose-up: Start MongoDB and gRPC server using docker compose
compose-up: vendor
	@echo "==> Starting services with Docker Compose..."
	docker compose up -d

## compose-down: Stop all docker compose services
compose-down:
	@echo "==> Stopping Docker Compose services..."
	docker compose down

## clean: Remove build artifacts and temporary files
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/ coverage.out coverage.html swaggerapi vendor/

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
