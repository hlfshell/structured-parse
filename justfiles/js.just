# JavaScript-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Build JavaScript package
build:
    @echo "🔨 Building JavaScript..."
    @if [ -f "{{project_root}}/js/package.json" ]; then \
        cd {{project_root}}/js && npm run build; \
    else \
        echo "❌ JavaScript build configuration not found (package.json)"; \
        exit 1; \
    fi
    @echo "✅ JavaScript build complete!"


# Publish JavaScript package to npm
publish:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Publishing JavaScript package..."
    
    # Check if we're on a git tag
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before publishing."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    echo "✅ Current git tag: $CURRENT_TAG"
    
    if [ -f "{{project_root}}/js/package.json" ]; then
        cd {{project_root}}/js
        npm publish --access public
    else
        echo "❌ JavaScript publish configuration not found"
        exit 1
    fi
    @echo "✅ JavaScript publish complete!"

# Show available JavaScript commands
help:
    @echo "Available JavaScript commands:"
    @echo "  js build    - Build JavaScript package"
    @echo "  js test     - Run JavaScript tests"
    @echo "  js publish  - Publish to npm (requires git tag)"

