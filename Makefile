.PHONY: help dev prod test clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Start development environment with hot-reload
	docker-compose -f docker-compose.override.yml up --build

prod: ## Start production environment
	docker-compose up --build

test: ## Run all tests
	go test ./... -v -cover

test-integration: ## Run integration tests only
	go test ./test/integration/... -v

clean: ## Clean up containers and volumes
	docker-compose down -v
	rm -rf tmp/
