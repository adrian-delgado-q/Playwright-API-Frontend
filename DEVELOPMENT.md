# Development Guide

## Getting Started

### Prerequisites
- Node.js 18+ 
- Go 1.21+
- Git

### Setup
1. Clone the repository
2. Install dependencies:
   ```bash
   npm install
   cd books_api && go mod tidy
   ```

3. Install Playwright browsers:
   ```bash
   npx playwright install --with-deps
   ```

## Development Workflow

### 1. Code Quality
Before committing any code, ensure it passes formatting and linting:

```bash
# Format code
npm run format

# Check for linting issues
npm run lint

# Fix auto-fixable linting issues
npm run lint:fix
```

### 2. Testing
Run tests locally before pushing:

```bash
# Start the API server (in one terminal)
npm run api:start

# Run all tests (in another terminal)
npm test

# Or run specific test suites
npm run test:api
npm run test:smoke
npm run test:crud
```

### 3. Building
Test that the application builds correctly:

```bash
npm run api:build
```

## CI/CD Pipeline

The GitHub Actions pipeline automatically runs on:
- Push to `main` or `develop` branches
- Pull requests to `main`

### Pipeline Stages

#### Stage 1: Format & Lint
- Prettier format checking
- ESLint code quality checks
- Results cached for faster subsequent runs

#### Stage 2: Test
- Go integration tests
- Playwright E2E tests (multiple test types in parallel)
- Test artifacts uploaded on failure

#### Stage 3: Build
- Go binary compilation
- Binary functionality verification
- Release artifacts created and uploaded

#### Stage 4: Deploy (Production only)
- Triggered only on `main` branch pushes
- Deployment simulation
- Post-deployment smoke tests
- Notification system

### Security Scanning
- npm audit for Node.js vulnerabilities
- Go dependency vulnerability checks

## File Structure

```
├── .github/workflows/ci.yml    # GitHub Actions CI/CD pipeline
├── .gitignore                  # Git ignore patterns
├── .prettierrc.json           # Prettier configuration
├── .prettierignore            # Prettier ignore patterns
├── eslint.config.mjs          # ESLint configuration (flat config)
├── scripts/pre-commit         # Pre-commit hook script
└── ... (existing files)
```

## Best Practices

### Code Quality
- All code must pass ESLint checks
- Use Prettier for consistent formatting
- Follow TypeScript best practices
- Write comprehensive tests

### Git Workflow
- Use feature branches for development
- Write clear commit messages
- Ensure CI passes before merging
- Use the pre-commit hook for local validation

### Testing
- Write tests for new features
- Maintain good test coverage
- Use appropriate test tags (@smoke, @api, @crud, etc.)
- Test both happy path and error scenarios

## Troubleshooting

### ESLint Issues
```bash
# Check specific file
npx eslint path/to/file.ts

# Fix auto-fixable issues
npx eslint path/to/file.ts --fix
```

### Prettier Issues
```bash
# Check specific file
npx prettier --check path/to/file.ts

# Format specific file
npx prettier --write path/to/file.ts
```

### CI Failures
1. Check the GitHub Actions tab for detailed logs
2. Run the same commands locally
3. Ensure all dependencies are up to date
4. Check for merge conflicts

## Environment Variables

For deployment and advanced configuration, you may need:

```bash
# API Configuration
export API_PORT=8080
export DB_PATH=./books.db

# Test Configuration  
export TEST_TIMEOUT=30000
export HEADLESS=true
```
