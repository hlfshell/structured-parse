# structured-parse - Task Runner
# Root Justfile using delegation-based, git-style subcommand pattern

# Import sub-justfiles as modules
mod go 'justfiles/go.just'
mod python 'justfiles/python.just'
mod ts 'justfiles/ts.just'
mod js 'justfiles/js.just'

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default:
    @echo "═══════════════════════════════════════════════════════════════"
    @echo "  🤖 structured-parse 🧠"
    @echo "═══════════════════════════════════════════════════════════════"
    @echo ""
    @echo "  🧹 clean    - Clean all build artifacts and copied files"
    @echo "  🔨 build    - Build all languages"
    @echo "  🧪 test     - Run tests for all languages"
    @echo "  📦 publish  - Publish to package managers (checks git tags)"
    @echo ""
    @echo "Language-specific commands:"
    @echo "  go          - Go-specific commands (build, test, bench)"
    @echo "  python      - Python-specific commands (build, test, publish)"
    @echo "  ts          - TypeScript-specific commands (build, test, publish)"
    @echo "  js          - JavaScript-specific commands (build, test, publish)"
    @echo ""
    @echo "═══════════════════════════════════════════════════════════════"

# Show this help message
help: default

# ============================================================================
# Clean all build artifacts and copied files
# ============================================================================
clean:
    @echo "🧹 Cleaning all build artifacts and copied files..."
    @echo ""
    @just python clean
    @echo ""
    @just ts clean
    @echo ""
    @just js clean
    @echo ""
    @echo "✅ All clean complete!"

# ============================================================================
# Build all languages
# ============================================================================
build:
    @echo "🔨 Building all languages..."
    @echo ""
    @just go build
    @echo ""
    @echo "🔨 Building WASM modules..."
    @just go build-wasm
    @echo ""
    @just python build
    @echo ""
    @just ts build
    @echo ""
    @just js build
    @echo ""
    @echo "✅ All builds complete!"

# ============================================================================
# Test all languages
# ============================================================================
test:
    @echo "🧪 Running tests for all languages..."
    @echo ""
    @just go test
    @echo ""
    @just python test
    @echo ""
    @just ts test
    @echo ""
    @echo "✅ All tests complete!"

# ============================================================================
# Publish to package managers (checks git tags first)
# ============================================================================
publish:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Publishing to package managers..."
    echo ""
    
    # Check if we're on a git tag
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before publishing."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    echo "✅ Current git tag: $CURRENT_TAG"
    echo ""
    
    # Check for uncommitted changes
    if ! git diff-index --quiet HEAD --; then
        echo "⚠️  Warning: You have uncommitted changes."
        read -p "Continue anyway? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    echo "📦 Publishing Python..."
    just python publish
    echo ""
    
    echo "📦 Publishing TypeScript..."
    just ts publish
    echo ""
    
    echo "📦 Publishing JavaScript..."
    just js publish
    echo ""
    
    echo "✅ Publish complete!"
