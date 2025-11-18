# TypeScript-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Install TypeScript dependencies
install:
    @echo "📦 Installing TypeScript dependencies..."
    @cd {{project_root}}/ts && npm install
    @echo "✅ TypeScript dependencies installed!"

# Build TypeScript package (requires WASM to be built first)
build:
    @echo "🔨 Building TypeScript..."
    @if [ -f "{{project_root}}/ts/package.json" ]; then \
        cd {{project_root}}/ts && npm run build; \
    else \
        echo "❌ TypeScript build configuration not found (package.json)"; \
        exit 1; \
    fi
    @echo "✅ TypeScript build complete!"

# Run TypeScript tests
test:
    @echo "🧪 Running TypeScript tests..."
    @if [ -f "{{project_root}}/ts/package.json" ]; then \
        cd {{project_root}}/ts && npm test; \
    else \
        echo "❌ TypeScript test configuration not found"; \
        exit 1; \
    fi
    @echo "✅ TypeScript tests complete!"

# Publish TypeScript package to npm
publish:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Publishing TypeScript package..."
    
    # Check if we're on a git tag
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before publishing."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    echo "✅ Current git tag: $CURRENT_TAG"
    
    if [ -f "{{project_root}}/ts/package.json" ]; then
        cd {{project_root}}/ts
        npm publish --access public
    else
        echo "❌ TypeScript publish configuration not found"
        exit 1
    fi
    @echo "✅ TypeScript publish complete!"

# Show available TypeScript commands
help:
    @echo "Available TypeScript commands:"
    @echo "  ts install  - Install dependencies"
    @echo "  ts build    - Build TypeScript package"
    @echo "  ts test     - Run TypeScript tests"
    @echo "  ts publish  - Publish to npm (requires git tag)"

