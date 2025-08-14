# Docker Development Guide

## Overview

This project has been containerized using Docker and Docker Compose for consistent development and deployment environments.

## Architecture

```
┌─────────────────┐    HTTP:3000     ┌─────────────────┐    HTTP:8080    ┌─────────────────┐
│   Web Browser   │ ◄─────────────► │   Nginx Frontend│ ◄─────────────► │   Go API Server │
│                 │                 │   (books-frontend) │               │   (books-api)   │
└─────────────────┘                 └─────────────────┘                 └─────────────────┘
                                              │                                    │
                                              └────── Proxy /api/* ──────────────┘
                                                                         │
                                                            ┌─────────────────┐
                                                            │ SQLite Database │
                                                            │   (books.db)    │
                                                            └─────────────────┘
```

## Services

### API Service (books-api)
- **Port**: 8080
- **Image**: Built from `books_api/Dockerfile`
- **Health Check**: `GET /health`
- **Database**: SQLite persisted in Docker volume

### Frontend Service (books-frontend) 
- **Port**: 3000 (mapped to nginx port 80)
- **Image**: Built from `frontend/Dockerfile`
- **Serves**: Static HTML/CSS/JS
- **Proxy**: Routes `/api/*` requests to API service
- **Health Check**: `GET /health`

## Quick Start

### 1. Setup and Validate
```bash
# Run the Docker setup script
npm run docker:setup
```

### 2. Development Workflow
```bash
# Start all services
npm run docker:up

# View logs
npm run docker:logs

# Run tests against Docker services
npm run docker:test

# Stop services
npm run docker:down
```

### 3. Access the Application
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080/api/v1/books
- **Health Checks**: 
  - API: http://localhost:8080/health
  - Frontend: http://localhost:3000/health

## Available Docker Commands

### Development
```bash
npm run docker:build      # Build all images
npm run docker:up         # Start services in background
npm run docker:down       # Stop and remove containers
npm run docker:logs       # View service logs
npm run docker:dev        # Start with development overrides
```

### Testing
```bash
npm run docker:test:build # Build test environment
npm run docker:test:up    # Start test services
npm run docker:test:down  # Stop test services
npm run docker:test       # Full test cycle (build, test, cleanup)
```

### Maintenance
```bash
npm run docker:clean      # Remove containers and prune system
npm run docker:setup      # Full setup and validation
```

## Docker Compose Files

### `docker-compose.yml`
Main production-ready configuration with:
- Health checks
- Volume persistence
- Network isolation
- Optimized for production

### `docker-compose.dev.yml`
Development overrides with:
- Debug mode for API
- Live file mounting
- Development-friendly settings

### `docker-compose.test.yml`
Testing configuration with:
- Faster health checks
- Isolated test environment
- Separate volumes

## Environment Variables

### API Service
- `GIN_MODE`: `release` (production) or `debug` (development)
- `DB_PATH`: Path to SQLite database file (default: `/data/books.db`)

### Frontend Service
- Uses nginx configuration for routing and proxying

## Volumes

### `api_data`
Persistent storage for the SQLite database ensuring data survives container restarts.

### `test_api_data` 
Isolated storage for test runs, automatically cleaned up after tests.

## Development Tips

### Debugging
```bash
# Check container status
docker-compose ps

# View specific service logs
docker-compose logs api
docker-compose logs frontend

# Execute commands in running containers
docker-compose exec api sh
docker-compose exec frontend sh

# Check health status
curl http://localhost:8080/health
curl http://localhost:3000/health
```

### Database Access
```bash
# Access SQLite database
docker-compose exec api sh
# Inside container:
sqlite3 /data/books.db
.tables
.quit
```

### Rebuilding Images
```bash
# Rebuild specific service
docker-compose build api
docker-compose build frontend

# Force rebuild without cache
docker-compose build --no-cache
```

## CI/CD Integration

The GitHub Actions pipeline has been updated to:

1. **Format & Lint** - Code quality checks
2. **Test** - Docker-based testing with service health checks
3. **Build** - Docker image building and validation
4. **Deploy** - Container-based deployment simulation

### Test Stage
- Builds test-specific Docker images
- Starts services with health checks
- Runs Playwright tests against containerized services
- Captures logs on failure
- Cleans up test environment

### Build Stage
- Builds production Docker images
- Validates image functionality
- Saves images as CI artifacts
- Creates release archives

## Troubleshooting

### Port Conflicts
If ports 3000 or 8080 are already in use:
```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080

# Stop conflicting services or modify docker-compose.yml ports
```

### Database Issues
```bash
# Reset database
docker-compose down -v  # This removes volumes
docker-compose up -d    # Fresh database with seed data
```

### Image Build Issues
```bash
# Clean Docker cache
docker system prune -a

# Rebuild from scratch
docker-compose build --no-cache
```

### Health Check Failures
```bash
# Check service logs
docker-compose logs api
docker-compose logs frontend

# Manual health check
curl -v http://localhost:8080/health
curl -v http://localhost:3000/health
```

## Production Deployment

For production deployment:

1. Use the main `docker-compose.yml`
2. Set appropriate environment variables
3. Configure external reverse proxy/load balancer
4. Set up proper monitoring and logging
5. Configure backup strategies for the database volume

### Example Production Setup
```bash
# Production environment
export GIN_MODE=release
export API_PORT=8080
export FRONTEND_PORT=3000

# Deploy
docker-compose up -d

# Monitor
docker-compose logs -f
```
