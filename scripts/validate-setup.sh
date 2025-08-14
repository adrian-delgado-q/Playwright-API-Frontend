#!/bin/bash

# Development Setup Validation Script
echo "🔧 Validating development setup..."

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: Run this script from the project root directory"
    exit 1
fi

# Check Node.js version
echo "📦 Checking Node.js version..."
node_version=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$node_version" -lt 18 ]; then
    echo "⚠️  Warning: Node.js 18+ is recommended (current: $(node --version))"
else
    echo "✅ Node.js version: $(node --version)"
fi

# Check Go version
echo "🐹 Checking Go version..."
if command -v go >/dev/null 2>&1; then
    echo "✅ Go version: $(go version)"
else
    echo "❌ Go is not installed"
    exit 1
fi

# Check if dependencies are installed
echo "📚 Checking dependencies..."
if [ ! -d "node_modules" ]; then
    echo "⚠️  Node.js dependencies not installed. Running npm install..."
    npm install
fi

# Check Go dependencies
echo "🔍 Checking Go dependencies..."
cd books_api
if [ ! -f "go.sum" ]; then
    echo "⚠️  Go dependencies not installed. Running go mod tidy..."
    go mod tidy
fi
cd ..

# Test formatting
echo "🎨 Testing Prettier formatting..."
if npm run format:check >/dev/null 2>&1; then
    echo "✅ Code formatting is correct"
else
    echo "⚠️  Code formatting issues found. Run 'npm run format' to fix"
fi

# Test linting
echo "🔍 Testing ESLint..."
if npm run lint >/dev/null 2>&1; then
    echo "✅ No linting issues found"
else
    echo "⚠️  Linting issues found. Run 'npm run lint:fix' to fix auto-fixable issues"
fi

# Test Go build
echo "🏗️  Testing Go build..."
cd books_api
if go build -o bin/books_api main.go >/dev/null 2>&1; then
    echo "✅ Go build successful"
    rm -f bin/books_api
else
    echo "❌ Go build failed"
    cd ..
    exit 1
fi
cd ..

# Check if Playwright is installed
echo "🎭 Checking Playwright installation..."
if npx playwright --version >/dev/null 2>&1; then
    echo "✅ Playwright is installed: $(npx playwright --version)"
else
    echo "⚠️  Playwright not found. Run 'npx playwright install --with-deps'"
fi

echo ""
echo "🎉 Development setup validation complete!"
echo ""
echo "📝 Next steps:"
echo "   1. Run 'npm run api:start' to start the API server"
echo "   2. Run 'npm test' to run all tests"
echo "   3. Check out DEVELOPMENT.md for detailed development guide"
echo ""
echo "🚀 Happy coding!"
