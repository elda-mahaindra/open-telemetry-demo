.PHONY: help up down

# Default target
help: ## Show available commands
	@echo "OpenTelemetry Demo Commands:"
	@echo ""
	@echo "  make up    - Start all services"
	@echo "  make down  - Stop all services"
	@echo ""

up: ## Start all services
	@echo "🚀 Starting services..."
	docker compose up -d
	@echo "✅ Services started!"
	@echo "📊 Jaeger UI: http://localhost:16686"
	@echo "🚀 API: http://localhost:4000"

down: ## Stop all services
	@echo "🛑 Stopping services..."
	docker compose down
	@echo "✅ Services stopped!" 