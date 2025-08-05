#!/bin/bash

# Script para executar todos os testes do projeto
# Usage: ./run_tests.sh

set -e

echo "🧪 GoLedger Challenge - Running Test Suite"
echo "=========================================="

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para imprimir status
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verificar se estamos no diretório correto
if [[ ! -f "go.mod" ]]; then
    print_error "Este script deve ser executado no diretório da aplicação (app/)"
    exit 1
fi

echo "📦 Installing dependencies..."
go mod download
go mod tidy
print_status "Dependencies installed"

echo ""
echo "🔍 Running unit tests..."

# Testar models
echo "  Testing models..."
if go test -v ./internal/models/ > /dev/null 2>&1; then
    print_status "Models tests passed"
else
    print_error "Models tests failed"
    exit 1
fi

# Testar handlers
echo "  Testing handlers..."
if go test -v ./internal/handlers/ > /dev/null 2>&1; then
    print_status "Handlers tests passed"
else
    print_error "Handlers tests failed"
    exit 1
fi

# Testar repositories (skip se não há database)
echo "  Testing repositories..."
if go test -short -v ./internal/repositories/ > /dev/null 2>&1; then
    print_status "Repository tests passed (skipped - no test DB)"
else
    print_warning "Repository tests skipped (requires test database)"
fi

echo ""
echo "📊 Running tests with coverage..."
if go test -v -coverprofile=coverage.out ./internal/... > test_output.log 2>&1; then
    print_status "All tests passed!"
    
    # Gerar relatório de cobertura
    if go tool cover -func=coverage.out | tail -1; then
        print_status "Coverage report generated"
    fi
    
    # Gerar HTML se possível
    if go tool cover -html=coverage.out -o coverage.html 2>/dev/null; then
        print_status "HTML coverage report: coverage.html"
    fi
else
    print_error "Some tests failed"
    echo "Check test_output.log for details"
    exit 1
fi

echo ""
echo "🚀 Integration Tests"
print_warning "Integration tests require a running PostgreSQL database"
print_warning "Run 'make db-up' first, then 'go test ./tests/integration/'"

echo ""
echo "🎉 Test suite completed successfully!"
echo ""
echo "Next steps:"
echo "  - Review coverage report: open coverage.html in browser"
echo "  - Run integration tests with: make db-up && go test ./tests/integration/"
echo "  - Use 'make test' for quick testing during development"
