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

# Clean TypeScript build artifacts and copied files
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧹 Cleaning TypeScript artifacts..."
    cd {{project_root}}/ts
    # Remove build artifacts
    rm -rf dist/ *.tsbuildinfo
    # Remove copied files
    rm -f README.md LICENSE
    echo "✅ TypeScript clean complete!"

# Build TypeScript package (requires WASM to be built first)
build: update
    @echo "🔨 Building TypeScript..."
    @if [ -f "{{project_root}}/ts/package.json" ]; then \
        if [ ! -f "{{project_root}}/ts/wasm/structured-parse.wasm" ]; then \
            echo "⚠️  WASM files not found. Building WASM first..."; \
            cd {{project_root}} && just go build-wasm-js; \
        fi; \
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

# Update TypeScript package: sync version, copy README and LICENSE
update: clean
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Get current git tag (strip 'v' prefix if present)
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before updating."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    if [ ! -f "{{project_root}}/ts/package.json" ]; then
        exit 0
    fi
    
    # Remove 'v' prefix if present for version comparison
    VERSION="${CURRENT_TAG#v}"
    
    # Check current version in package.json
    CURRENT_VERSION=$(cd {{project_root}}/ts && node -p "require('./package.json').version" 2>/dev/null || echo "")
    
    UPDATED=false
    
    # Update version if needed
    if [ "$CURRENT_VERSION" != "$VERSION" ]; then
        cd {{project_root}}/ts
        # Use node to update package.json version
        node -e "const fs = require('fs'); const pkg = JSON.parse(fs.readFileSync('package.json', 'utf8')); pkg.version = '$VERSION'; fs.writeFileSync('package.json', JSON.stringify(pkg, null, 2) + '\n');"
        UPDATED=true
    fi
    
    # Ensure README.md and LICENSE are in files array
    cd {{project_root}}/ts
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
    if [ ! -f "{{project_root}}/ts/README.md" ] || ! cmp -s "{{project_root}}/ts_js.md" "{{project_root}}/ts/README.md" 2>/dev/null; then
        cp "{{project_root}}/ts_js.md" "{{project_root}}/ts/README.md"
        UPDATED=true
    fi
    
    # Copy LICENSE if it doesn't exist
    if [ ! -f "{{project_root}}/ts/LICENSE" ]; then
        cp "{{project_root}}/LICENSE" "{{project_root}}/ts/LICENSE"
        UPDATED=true
    fi
    
    exit 0

# Publish TypeScript package to npm
publish: build
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Run update first (will fail if not on a tag)
    just ts update
    
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
    echo "✅ TypeScript publish complete!"

# Show available TypeScript commands
help:
    @echo "Available TypeScript commands:"
    @echo "  ts install  - Install dependencies"
    @echo "  ts clean    - Remove build artifacts and copied files"
    @echo "  ts build    - Build TypeScript package"
    @echo "  ts test     - Run TypeScript tests"
    @echo "  ts update   - Update version, README, and LICENSE"
    @echo "  ts publish  - Publish to npm (requires git tag, runs update first)"

