#!/bin/bash

# Docker Setup and Validation Script
echo "🐳 Setting up and validating Docker environment..."

# Check if Docker is installed
if ! command -v docker >/dev/null 2>&1; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker >/dev/null 2>&1; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is available (try both v1 and v2 syntax)
if command -v docker-compose >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker compose"
else
    echo "❌ Docker Compose is not available. Please install Docker Compose."
    exit 1
fi

# Check if Docker daemon is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker daemon is not running. Please start Docker first."
    exit 1
fi

echo "✅ Docker: $(docker --version)"
echo "✅ Docker Compose: $($DOCKER_COMPOSE version)"

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Run this script from the project root directory"
    exit 1
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
$DOCKER_COMPOSE down -v >/dev/null 2>&1 || true
$DOCKER_COMPOSE -f docker-compose.test.yml down -v >/dev/null 2>&1 || true

# Build the images
echo "🏗️  Building Docker images..."
if $DOCKER_COMPOSE build; then
    echo "✅ Docker images built successfully"
else
    echo "❌ Failed to build Docker images"
    exit 1
fi

# Start the services
echo "🚀 Starting services..."
if $DOCKER_COMPOSE up -d; then
    echo "✅ Services started successfully"
else
    echo "❌ Failed to start services"
    exit 1
fi

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1 && curl -f http://localhost:3000/health >/dev/null 2>&1; then
        echo "✅ All services are healthy!"
        break
    fi
    echo "Waiting for services... ($i/30)"
    sleep 2
done

# Test the services
echo "🧪 Testing services..."

# Test API
echo "Testing API endpoints..."
if curl -f http://localhost:8080/health >/dev/null 2>&1; then
    echo "✅ API health check passed"
else
    echo "❌ API health check failed"
fi

if curl -f http://localhost:8080/api/v1/books >/dev/null 2>&1; then
    echo "✅ API books endpoint accessible"
else
    echo "❌ API books endpoint failed"
fi

# Test Frontend
echo "Testing Frontend..."
if curl -f http://localhost:3000/health >/dev/null 2>&1; then
    echo "✅ Frontend health check passed"
else
    echo "❌ Frontend health check failed"
fi

if curl -f http://localhost:3000/ >/dev/null 2>&1; then
    echo "✅ Frontend accessible"
else
    echo "❌ Frontend failed"
fi

# Show running containers
echo ""
echo "📋 Running containers:"
$DOCKER_COMPOSE ps

echo ""
echo "🎉 Docker setup validation complete!"
echo ""
echo "📝 Available commands:"
echo "   npm run docker:up      - Start services"
echo "   npm run docker:down    - Stop services"
echo "   npm run docker:logs    - View logs"
echo "   npm run docker:test    - Run tests against Docker"
echo ""
echo "🌐 Access the application:"
echo "   Frontend: http://localhost:3000"
echo "   API:      http://localhost:8080/api/v1/books"
echo "   Health:   http://localhost:8080/health"
echo ""
echo "To stop the services, run: npm run docker:down"
