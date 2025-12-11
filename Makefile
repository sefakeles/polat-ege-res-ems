.PHONY: build run test clean docker lint fmt tidy help

# Build variables
BINARY_NAME=ems
MAIN_PATH=./cmd/ems
DOCKER_IMAGE=solservices-gyongyoshalasz-ems
VERSION?=dev

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go mod tidy
	CGO_ENABLED=1 go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

run: build ## Build and run the application
	./$(BINARY_NAME)

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests
	go test -v -tags=integration ./test/integration/...

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

clean: ## Clean build artifacts
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

docker: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) -f build/Dockerfile .

docker-run: docker ## Run Docker container
	docker run -d --name $(BINARY_NAME) \
		-p 8080:8080 \
		-v $(PWD)/configs:/app/configs:ro \
		-v $(PWD)/data:/app/data \
		$(DOCKER_IMAGE):$(VERSION)

docker-stop: ## Stop and remove Docker container
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

tidy: ## Tidy dependencies
	go mod tidy
	go mod verify

deps: ## Install development dependencies
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev: ## Run with live reload (requires air)
	air

install: build ## Install the binary
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

service-install: install ## Install as systemd service
	sudo ./scripts/install.sh

service-uninstall: ## Uninstall systemd service
	sudo systemctl stop solservices-gyongyoshalasz-ems || true
	sudo systemctl disable solservices-gyongyoshalasz-ems || true
	sudo rm -f /etc/systemd/system/solservices-gyongyoshalasz-ems.service
	sudo systemctl daemon-reload

backup: ## Create backup
	./scripts/backup.sh

monitor: ## Run monitoring check
	./scripts/monitor.sh

update: ## Update to new version
	./scripts/update.sh -v $(VERSION) -b -r

# Development targets
dev-deps: deps ## Install all development dependencies
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest

generate: ## Generate code (mocks, swagger, etc.)
	go generate ./...

security: ## Run security checks
	gosec ./...

mod-check: ## Check for module updates
	go list -u -m all

# CI/CD targets
ci-test: ## Run CI tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

ci-build: ## Build for CI
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
