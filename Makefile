# Claude Code Parser Makefile

.PHONY: help setup build run generate clean test fmt vet

# Default target
help:
	@echo "Claude Code Parser - Available targets:"
	@echo "  setup     - Initialize Go module and install dependencies"
	@echo "  generate  - Generate Go code from templ templates"
	@echo "  build     - Build the application"
	@echo "  run       - Run with default JSONL file"
	@echo "  run-open  - Run with default file and open HTML in browser"
	@echo "  run-file  - Run with custom JSONL file (make run-file FILE=your-file.jsonl)"
	@echo "  server    - Start web server on port 8080"
	@echo "  server-port - Start web server on custom port (make server-port PORT=3000)"
	@echo "  fmt       - Format Go code"
	@echo "  vet       - Run go vet"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean generated files and build artifacts"
	@echo ""
	@echo "Example usage:"
	@echo "  make setup"
	@echo "  make generate"
	@echo "  make run"
	@echo "  make run-file FILE=data.jsonl OUTPUT=output.html"

# Setup project dependencies
setup:
	@echo "Setting up Go module and dependencies..."
	go mod init claude-code-parser 2>/dev/null || true
	go get -tool github.com/a-h/templ/cmd/templ@latest
	go get github.com/a-h/templ
	go mod tidy
	@echo "✅ Setup complete"

# Generate Go code from templ templates
generate:
	@echo "Generating Go code from templ templates..."
	go tool templ generate
	@echo "✅ Templates generated"

# Build the application
build: generate
	@echo "Building application..."
	go build -o claude-code-parser .
	@echo "✅ Build complete"

# Default run with the existing JSONL file
run: generate
	@echo "Running with default JSONL file..."
	go run . 91cc2e2a-2d04-46ba-a5cf-5fcadf00f1da.jsonl conversation.html
	@echo "✅ HTML output written to conversation.html"

# Run and open the HTML file
run-open: run
	@echo "Opening HTML file..."
	open conversation.html || xdg-open conversation.html || echo "Please open conversation.html manually"

# Run with custom file
FILE ?= 91cc2e2a-2d04-46ba-a5cf-5fcadf00f1da.jsonl
OUTPUT ?= conversation.html
run-file: generate
	@echo "Running with file: $(FILE) -> $(OUTPUT)"
	go run . $(FILE) $(OUTPUT)
	@echo "✅ HTML output written to $(OUTPUT)"

# Start web server
server: generate
	@echo "Starting web server on port 8080..."
	go run . --server

# Start web server with custom port
PORT ?= 8080
server-port: generate
	@echo "Starting web server on port $(PORT)..."
	go run . --server --port $(PORT)

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ Vet passed"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "✅ Tests passed"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f claude-code-parser
	rm -f templates_templ.go
	rm -f *.html
	@echo "✅ Clean complete"

# Development workflow - setup, generate, and run
dev: setup generate run

# Full check - format, vet, test, build
check: fmt vet test build
	@echo "✅ All checks passed"

# Show project structure
structure:
	@echo "Project structure:"
	@tree -I 'go.sum|*.html' . 2>/dev/null || find . -type f -name "*.go" -o -name "*.templ" -o -name "*.jsonl" -o -name "Makefile" | sort