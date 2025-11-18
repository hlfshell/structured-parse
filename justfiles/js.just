# JavaScript-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Build JavaScript package (requires WASM to be built first)
build:
    @echo "🔨 Building JavaScript..."
    @if [ -f "{{project_root}}/js/package.json" ]; then \
        if [ ! -f "{{project_root}}/js/wasm/structured-parse.wasm" ]; then \
            echo "⚠️  WASM files not found. Building WASM first..."; \
            cd {{project_root}} && just go build-wasm-js; \
        fi; \
        cd {{project_root}}/js && npm run build; \
    else \
        echo "❌ JavaScript build configuration not found (package.json)"; \
        exit 1; \
    fi
    @echo "✅ JavaScript build complete!"

# Update JavaScript package: sync version, copy README and LICENSE
update:
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Get current git tag (strip 'v' prefix if present)
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before updating."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    if [ ! -f "{{project_root}}/js/package.json" ]; then
        exit 0
    fi
    
    # Remove 'v' prefix if present for version comparison
    VERSION="${CURRENT_TAG#v}"
    
    # Check current version in package.json
    CURRENT_VERSION=$(cd {{project_root}}/js && node -p "require('./package.json').version" 2>/dev/null || echo "")
    
    UPDATED=false
    
    # Update version if needed
    if [ "$CURRENT_VERSION" != "$VERSION" ]; then
        cd {{project_root}}/js
        # Use node to update package.json version
        node -e "const fs = require('fs'); const pkg = JSON.parse(fs.readFileSync('package.json', 'utf8')); pkg.version = '$VERSION'; fs.writeFileSync('package.json', JSON.stringify(pkg, null, 2) + '\n');"
        UPDATED=true
    fi
    
    # Ensure README.md and LICENSE are in files array
    cd {{project_root}}/js
    node -e "
        const fs = require('fs');
        const pkg = JSON.parse(fs.readFileSync('package.json', 'utf8'));
        if (!pkg.files) pkg.files = [];
        let changed = false;
        if (!pkg.files.includes('README.md')) {
            pkg.files.push('README.md');
            changed = true;
        }
        if (!pkg.files.includes('LICENSE')) {
            pkg.files.push('LICENSE');
            changed = true;
        }
        if (changed) {
            fs.writeFileSync('package.json', JSON.stringify(pkg, null, 2) + '\n');
        }
    " 2>/dev/null || true
    
    # Copy README if it doesn't exist or is different
    if [ ! -f "{{project_root}}/js/README.md" ] || ! cmp -s "{{project_root}}/ts_js.md" "{{project_root}}/js/README.md" 2>/dev/null; then
        cp "{{project_root}}/ts_js.md" "{{project_root}}/js/README.md"
        UPDATED=true
    fi
    
    # Copy LICENSE if it doesn't exist
    if [ ! -f "{{project_root}}/js/LICENSE" ]; then
        cp "{{project_root}}/LICENSE" "{{project_root}}/js/LICENSE"
        UPDATED=true
    fi
    
    exit 0

# Publish JavaScript package to npm
publish:
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Run update first (will fail if not on a tag)
    just js update
    
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
    echo "✅ JavaScript publish complete!"

# Show available JavaScript commands
help:
    @echo "Available JavaScript commands:"
    @echo "  js build    - Build JavaScript package"
    @echo "  js test     - Run JavaScript tests"
    @echo "  js update   - Update version, README, and LICENSE"
    @echo "  js publish  - Publish to npm (requires git tag, runs update first)"

