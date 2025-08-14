# Books API Project

This project demonstrates a complete CRUD API for managing books, built with Go and tested with Playwright.

## Project Structure

```
├── books_api/              # Go API backend
│   ├── main.go            # Main API server
│   ├── main_test.go       # Integration tests
│   ├── go.mod             # Go module dependencies
│   ├── go.sum             # Go module checksums
│   ├── Makefile           # Build and run commands
│   └── books.db           # SQLite database
├── tests/                  # Playwright E2E tests
│   ├── books-api.spec.ts  # API endpoint tests
│   ├── books-web.spec.ts  # Web interface tests
│   └── examples/          # Example tests
├── books_demo.html         # Demo web interface
├── playwright.config.ts   # Playwright configuration
└── package.json           # Node.js dependencies and scripts
```

## Features

### Books API (Go + Gorilla Mux + GORM + SQLite)

- **GET** `/api/v1/books` - List all books
- **GET** `/api/v1/books/{id}` - Get book by ID
- **POST** `/api/v1/books` - Create new book
- **PUT** `/api/v1/books/{id}` - Update book
- **DELETE** `/api/v1/books/{id}` - Delete book
- **GET** `/health` - Health check endpoint

### Features:

- ✅ Full CRUD operations
- ✅ SQLite database with GORM ORM
- ✅ CORS support for web interface
- ✅ Database seeding with sample data
- ✅ Input validation
- ✅ Error handling
- ✅ RESTful API design

### Testing

#### Integration Tests (Go)

- Unit tests for all API endpoints
- Database mocking with in-memory SQLite
- Test coverage for CRUD operations

#### E2E Tests (Playwright)

- **API Tests**: Direct HTTP requests to test endpoints
- **Web Tests**: Browser automation testing the demo interface
- **Test Tags**: Organized with tags (@smoke, @api, @crud, @web, etc.)

### Demo Web Interface

- Modern, responsive HTML/CSS/JavaScript interface
- Real-time interaction with the API
- Add, edit, delete books functionality
- Beautiful gradient design
- Mobile-responsive layout

## Quick Start

### 1. Start the API Server

```bash
# Option 1: Using npm script
npm run api:start

# Option 2: Direct Go command
cd books_api && go run main.go
```

The API will be available at `http://localhost:8080`

### 2. Run Integration Tests

```bash
# Go integration tests
npm run api:test

# Or directly:
cd books_api && go test -v
```

### 3. Run Playwright E2E Tests

```bash
# All tests
npm test

# API tests only
npm run test:api

# Smoke tests
npm run test:smoke

# Web interface tests
npm run test:web
```

### 4. View Demo Interface

```bash
# Start HTTP server for demo
npm run serve:demo

# Then open: http://localhost:3000/books_demo.html
```

## Available Test Tags

- `@smoke` - Critical functionality tests
- `@api` - API endpoint tests
- `@web` - Web interface tests
- `@crud` - Create, Read, Update, Delete operations
- `@books` - Book-related functionality
- `@validation` - Input validation tests
- `@cors` - Cross-origin request tests
- `@integration` - Data consistency tests
- `@performance` - Concurrent request tests

## Sample API Usage

### Create a Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Architecture",
    "author": "Robert C. Martin",
    "isbn": "9780134494166",
    "year": 2017
  }'
```

### Get All Books

```bash
curl http://localhost:8080/api/v1/books
```

### Update a Book

```bash
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title"}'
```

### Delete a Book

```bash
curl -X DELETE http://localhost:8080/api/v1/books/1
```

## Database Schema

```sql
CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn TEXT UNIQUE NOT NULL,
    year INTEGER
);
```

## Dependencies

### Go Dependencies

- `github.com/gorilla/mux` - HTTP router
- `gorm.io/gorm` - ORM library
- `gorm.io/driver/sqlite` - SQLite driver

### Node.js Dependencies

- `@playwright/test` - E2E testing framework
- `@types/node` - TypeScript definitions

## Development

### Running Tests

```bash
# Run all Playwright tests
npm test

# Run with specific tags
npm run test:smoke
npm run test:api
npm run test:crud

# Run in headed mode (see browser)
npm run test:headed

# Run with UI mode
npm run test:ui
```

### Building

```bash
# Build API binary
npm run api:build

# Run Go tests
npm run api:test
```

### Makefile Commands (in books_api/)

```bash
make build      # Build binary
make run        # Run server
make test       # Run tests
make clean      # Clean build artifacts
make dev        # Run in development mode
```

## Development Tools

### Code Quality & Formatting

This project uses ESLint and Prettier to maintain code quality and consistent formatting.

#### Available Scripts

```bash
# Format all files with Prettier
npm run format

# Check if files are properly formatted
npm run format:check

# Run ESLint to check for code quality issues
npm run lint

# Automatically fix ESLint issues where possible
npm run lint:fix

# Run all format and lint checks (used in CI)
npm run ci:format
```

#### Pre-commit Hooks

You can optionally set up a pre-commit hook to automatically run formatting and linting:

```bash
# Link the pre-commit script
ln -sf ../../scripts/pre-commit .git/hooks/pre-commit
```

### CI/CD Pipeline

The project includes a comprehensive GitHub Actions CI/CD pipeline with 4 stages:

1. **Format & Lint** - Code quality checks with Prettier and ESLint
2. **Test** - Integration tests (Go) and E2E tests (Playwright)
3. **Build** - Compile the Go API binary and create release artifacts
4. **Deploy** - Deploy to production (currently simulated)

The pipeline includes:

- ✅ Multi-matrix testing (api, smoke, crud tests)
- ✅ Security vulnerability scanning
- ✅ Artifact uploads for builds and test reports
- ✅ Deployment simulation with smoke tests
- ✅ Parallel job execution for faster builds

## Architecture

```
┌─────────────────┐    HTTP/JSON    ┌─────────────────┐
│   Web Browser   │ ◄─────────────► │   Go API Server │
│  (books_demo.   │                 │    (port 8080)  │
│     html)       │                 │                 │
└─────────────────┘                 └─────────────────┘
                                               │
┌─────────────────┐                           │
│   Playwright    │                           ▼
│   Test Runner   │                 ┌─────────────────┐
│                 │                 │ SQLite Database │
└─────────────────┘                 │   (books.db)    │
                                    └─────────────────┘
```

This project demonstrates modern API development practices with comprehensive testing strategies using both integration testing and end-to-end browser automation.
