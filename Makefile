# Makefile for maintaining the Flux CLI plugin catalog.

.PHONY: build validate clean

build: ## Build the validate binary
	@cd cmd/validate && go fmt ./... && go mod tidy && go build -o ../../bin/validate .

validate: build ## Validate catalog and plugin manifests against JSON schemas
	@./bin/validate schemas/ catalog.yaml plugins/*.yaml

clean: ## Remove build artifacts
	@rm -rf bin/

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-20s %s\n", $$1, $$2}'
