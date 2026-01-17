#!/bin/bash
# scripts/dev-setup.sh

set -e

echo "=== Setting up development environment ==="

# Create config directory if it doesn't exist
mkdir -p config

# Create config file for local development (non-Docker)
cat > config/config.yaml <<EOF
server:
  port: 8080
  environment: development
  read_timeout: 15s
  write_timeout: 15s

database:
  host: localhost
  port: 5432
  user: postgres
  password: 
  dbname: invoices
  sslmode: disable
  max_conns: 25
  min_conns: 5

auth:
  jwt_secret: your-secret-key-change-in-production
  token_duration: 24h

jira:
  base_url: ""
  enabled: false

square:
  access_token: ""
  environment: sandbox
  enabled: false

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
EOF

echo "✓ Config file created"

# Stop and remove existing containers and volumes
echo ""
echo "Cleaning up existing containers and volumes..."
docker-compose down -v 2>/dev/null || true

# Rebuild images
echo ""
echo "Building Docker images..."
docker-compose build --no-cache

# Start containers (migrations will run automatically)
echo ""
echo "Starting Docker containers..."
docker-compose up -d

# Wait for services to be healthy
echo ""
echo "Waiting for services to be ready..."
sleep 8

# Check if containers are running
echo ""
echo "Checking container status..."
docker-compose ps

# Check API logs
echo ""
echo "=== API Logs (last 20 lines) ==="
docker-compose logs --tail=20 api

echo ""
echo "=== Setup complete! ==="
echo ""
echo "Available commands:"
echo "  docker-compose logs -f api      # Follow API logs"
echo "  docker-compose logs -f postgres # Follow database logs"
echo "  docker-compose down             # Stop all services"
echo "  docker-compose up -d            # Start services"
echo ""
echo "API should be running at: http://localhost:8080"
echo "Health check: curl http://localhost:8080/api/v1/health"