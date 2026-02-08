.PHONY: help build run dev clean docker-build docker-run docker-stop frontend backend

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: frontend backend ## Build both frontend and backend

frontend: ## Build frontend
	cd frontend && npm install && npm run build

backend: ## Build backend
	mkdir -p backend/frontend
	cp -r frontend/dist backend/frontend/
	cd backend && go mod download && go build -o plocate-ui

dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

dev-backend: ## Run backend in development mode
	cd backend && go run main.go

docker-build: ## Build Docker image
	docker build -t plocate-ui:latest .

docker-run: ## Run Docker container
	docker compose up -d

docker-stop: ## Stop Docker container
	docker compose down

docker-logs: ## Show Docker logs
	docker compose logs -f

clean: ## Clean build artifacts
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/frontend
	rm -f backend/plocate-ui

test-backend: ## Run backend tests
	cd backend && go test ./...

install-deps: ## Install all dependencies
	cd frontend && npm install
	cd backend && go mod download
