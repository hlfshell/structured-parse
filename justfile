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
# Build all languages
# ============================================================================
build:
    @echo "🔨 Building all languages..."
    @echo ""
    @just go build
    @echo ""
    @if [ -f "{{project_root}}/python/pyproject.toml" ] || [ -f "{{project_root}}/python/setup.py" ]; then \
        just python build; \
    else \
        echo "⚠️  Python build not configured yet"; \
    fi
    @echo ""
    @if [ -f "{{project_root}}/ts/package.json" ]; then \
        just ts build; \
    else \
        echo "⚠️  TypeScript build not configured yet"; \
    fi
    @echo ""
    @if [ -f "{{project_root}}/js/package.json" ]; then \
        just js build; \
    else \
        echo "⚠️  JavaScript build not configured yet"; \
    fi
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
    @if [ -f "{{project_root}}/python/pyproject.toml" ] || [ -f "{{project_root}}/python/setup.py" ]; then \
        just python test; \
    else \
        echo "⚠️  Python tests not configured yet"; \
    fi
    @echo ""
    @if [ -f "{{project_root}}/ts/package.json" ]; then \
        just ts test; \
    else \
        echo "⚠️  TypeScript tests not configured yet"; \
    fi
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
    if [ -f "{{project_root}}/python/pyproject.toml" ] || [ -f "{{project_root}}/python/setup.py" ]; then
        just python publish
    else
        echo "⚠️  Python publish not configured yet"
    fi
    echo ""
    
    echo "📦 Publishing TypeScript..."
    if [ -f "{{project_root}}/ts/package.json" ]; then
        just ts publish
    else
        echo "⚠️  TypeScript publish not configured yet"
    fi
    echo ""
    
    echo "📦 Publishing JavaScript..."
    if [ -f "{{project_root}}/js/package.json" ]; then
        just js publish
    else
        echo "⚠️  JavaScript publish not configured yet"
    fi
    echo ""
    
    echo "✅ Publish complete!"
