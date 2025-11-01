# JEDI API Makefile
# Convenient commands for building, running, and managing the JEDI Docker container

.PHONY: help build run stop restart logs clean test test-encrypt test-decrypt status

# Default target
help:
	@echo "🚀 JEDI API - Docker Management Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make build          Build the JEDI Docker image"
	@echo "  make run            Run the JEDI container"
	@echo "  make stop           Stop the JEDI container"
	@echo "  make restart        Restart the JEDI container"
	@echo "  make logs           Show container logs"
	@echo "  make clean          Remove container and image"
	@echo "  make test           Run basic API tests"
	@echo "  make test-encrypt   Test encryption endpoint"
	@echo "  make test-decrypt   Test decryption endpoint"
	@echo "  make status         Show container status"
	@echo "  make help           Show this help message"
	@echo ""
	@echo "Quick start: make build && make run && make test"

# Build the Docker image
build:
	@echo "🔨 Building JEDI Docker image..."
	docker build -t jedi-encrypted .
	@echo "✅ Build completed!"

# Run the container
run:
	@echo "🚀 Starting JEDI container..."
	docker run -d -p 8081:8080 --name jedi-api jedi-encrypted
	@echo "✅ Container started on http://localhost:8081"
	@sleep 2
	@make status

# Stop the container
stop:
	@echo "🛑 Stopping JEDI container..."
	docker stop jedi-api 2>/dev/null || echo "Container not running"
	@echo "✅ Container stopped"

# Restart the container
restart:
	@echo "🔄 Restarting JEDI container..."
	docker restart jedi-api 2>/dev/null || (echo "Starting new container..." && make run)
	@echo "✅ Container restarted"
	@sleep 2
	@make status

# Show container logs
logs:
	@echo "📋 JEDI Container Logs:"
	@echo "----------------------------------------"
	docker logs jedi-api 2>/dev/null || echo "Container not found"

# Follow logs in real-time
logs-follow:
	@echo "📋 Following JEDI Container Logs (Ctrl+C to exit):"
	@echo "----------------------------------------"
	docker logs -f jedi-api 2>/dev/null || echo "Container not found"

# Clean up containers and images
clean:
	@echo "🧹 Cleaning up..."
	-docker stop jedi-api 2>/dev/null
	-docker rm jedi-api 2>/dev/null
	-docker rmi jedi-encrypted 2>/dev/null
	@echo "✅ Cleanup completed"

# Show container status
status:
	@echo "📊 Container Status:"
	@echo "----------------------------------------"
	@if docker ps | grep -q jedi-api; then \
		echo "✅ JEDI container is running"; \
		echo "🌐 URL: http://localhost:8081"; \
		echo "📊 Health: http://localhost:8081/health"; \
	else \
		echo "❌ JEDI container is not running"; \
	fi
	@echo ""

# Run basic API tests
test: test-health test-encrypt test-decrypt
	@echo "🎉 All tests completed successfully!"

# Test health endpoint
test-health:
	@echo "🏥 Testing health endpoint..."
	@curl -s http://localhost:8081/health | jq -r '.status' | grep -q "healthy" && echo "✅ Health check passed" || echo "❌ Health check failed"

# Test encryption endpoint
test-encrypt:
	@echo "🔐 Testing encryption endpoint..."
	@curl -s -X POST http://localhost:8081/encrypt \
		-H "Content-Type: application/json" \
		-d '{"uri":"test","message":"Hello JEDI!"}' \
		| jq -r '.success' | grep -q "true" && echo "✅ Encryption test passed" || echo "❌ Encryption test failed"

# Test decryption endpoint
test-decrypt:
	@echo "🔓 Testing decryption endpoint..."
	@ENCRYPTED=$$(curl -s -X POST http://localhost:8081/encrypt \
		-H "Content-Type: application/json" \
		-d '{"uri":"test","message":"Hello JEDI!"}' | jq -r '.data'); \
	curl -s -X POST http://localhost:8081/decrypt \
		-H "Content-Type: application/json" \
		-d "{\"uri\":\"test\",\"encryptedMessage\":\"$$ENCRYPTED\"}" \
		| jq -r '.success' | grep -q "true" && echo "✅ Decryption test passed" || echo "❌ Decryption test failed"

# Full encryption/decryption cycle test
test-cycle:
	@echo "🔄 Testing full encryption/decryption cycle..."
	@MESSAGE="Test message at $$(date +%s)"
	@URI="test/cycle"
	@echo "Original: $$MESSAGE"
	@ENCRYPTED=$$(curl -s -X POST http://localhost:8081/encrypt \
		-H "Content-Type: application/json" \
		-d "{\"uri\":\"$$URI\",\"message\":\"$$MESSAGE\"}" | jq -r '.data'); \
	DECRYPTED=$$(curl -s -X POST http://localhost:8081/decrypt \
		-H "Content-Type: application/json" \
		-d "{\"uri\":\"$$URI\",\"encryptedMessage\":\"$$ENCRYPTED\"}" | jq -r '.data'); \
	if [ "$$MESSAGE" = "$$DECRYPTED" ]; then \
		echo "✅ Full cycle test passed"; \
		echo "Original: $$MESSAGE"; \
		echo "Decrypted: $$DECRYPTED"; \
	else \
		echo "❌ Full cycle test failed"; \
		echo "Original: $$MESSAGE"; \
		echo "Decrypted: $$DECRYPTED"; \
	fi

# Performance test
perf-test:
	@echo "⚡ Running performance test..."
	@echo "Testing 10 encryption operations..."
	@time (for i in {1..10}; do \
		curl -s -X POST http://localhost:8081/encrypt \
			-H "Content-Type: application/json" \
			-d '{"uri":"perf/test","message":"Performance test message $$i"}' > /dev/null; \
	done)
	@echo "✅ Performance test completed"

# Show container statistics
stats:
	@echo "📈 Container Statistics:"
	@echo "----------------------------------------"
	@docker stats --no-stream jedi-api 2>/dev/null || echo "Container not running"

# Enter container shell
shell:
	@echo "🐚 Entering container shell..."
	@docker exec -it jedi-api /bin/bash 2>/dev/null || echo "Container not running"

# Development mode - run without Docker
dev:
	@echo "🛠️ Running in development mode (without Docker)..."
	@echo "Starting JEDI API on http://localhost:8080..."
	go run enhanced_main.go

# Production deployment with health check
deploy-prod:
	@echo "🚀 Deploying to production..."
	docker run -d \
		-p 8081:8080 \
		--name jedi-api-prod \
		--restart unless-stopped \
		--memory=512m \
		--cpus=0.5 \
		jedi-encrypted
	@echo "✅ Production deployment completed"
	@sleep 3
	@make status

# Backup container data (if any volumes are used)
backup:
	@echo "💾 Creating backup..."
	@mkdir -p backups
	@docker export jedi-api > backups/jedi-api-$(shell date +%Y%m%d-%H%M%S).tar 2>/dev/null || echo "No container to backup"
	@echo "✅ Backup completed"

# Monitor API in real-time
monitor:
	@echo "👀 Monitoring API (Ctrl+C to exit)..."
	@watch -n 2 'curl -s http://localhost:8081/health | jq ".status, .timestamp"'

# Show API endpoints
endpoints:
	@echo "📡 Available API Endpoints:"
	@echo "----------------------------------------"
	@echo "GET  http://localhost:8081/        - API Information"
	@echo "GET  http://localhost:8081/health   - Health Check"
	@echo "POST http://localhost:8081/encrypt  - Encrypt Message"
	@echo "POST http://localhost:8081/decrypt  - Decrypt Message"
	@echo ""

# Quick demo
demo:
	@echo "🎬 JEDI API Demo"
	@echo "=================="
	@echo ""
	@echo "1. Checking API status..."
	@curl -s http://localhost:8081/ | jq -r '.message'
	@echo ""
	@echo "2. Encrypting a secret message..."
	@RESULT=$$(curl -s -X POST http://localhost:8081/encrypt \
		-H "Content-Type: application/json" \
		-d '{"uri":"demo/secret","message":"This is a secret message for the demo!"}'); \
	echo "$$RESULT" | jq -r '"Encrypted: " + .data'; \
	echo "$$RESULT" | jq -r '"Performance: " + (.executionTime|tostring) + "μs, Energy: " + (.energyConsumptionJoules|tostring) + "J"';
	@echo ""
	@echo "3. Decrypting the message..."
	@ENCRYPTED=$$(echo "$$RESULT" | jq -r '.data'); \
	curl -s -X POST http://localhost:8081/decrypt \
		-H "Content-Type: application/json" \
		-d "{\"uri\":\"demo/secret\",\"encryptedMessage\":\"$$ENCRYPTED\"}" | jq -r '"Decrypted: " + .data';
	@echo ""
	@echo "✅ Demo completed successfully!"

# Install dependencies (for local development)
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@echo "✅ Dependencies installed"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Security scan
security-scan:
	@echo "🔒 Running security scan..."
	@echo "Checking for common vulnerabilities..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy image jedi-encrypted || echo "Trivy not available"
	@echo "✅ Security scan completed"