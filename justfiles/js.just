# JavaScript-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Clean JavaScript build artifacts and copied files
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧹 Cleaning JavaScript artifacts..."
    cd {{project_root}}/js
    # Remove build artifacts
    rm -rf dist/ package.json
    # Remove copied files
    rm -f README.md LICENSE
    echo "✅ JavaScript clean complete!"


# Build JavaScript package (requires TypeScript to be built first)
build: update
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "🔨 Building JavaScript package..."
    
    # Ensure TypeScript is built first
    if [ ! -d "{{project_root}}/ts/dist" ] || [ ! -f "{{project_root}}/ts/dist/index.js" ]; then
        echo "⚠️  TypeScript dist not found. Building TypeScript first..."
        cd {{project_root}} && just ts build
    fi
    
    # Create js/dist directory
    mkdir -p {{project_root}}/js/dist/wasm
    
    # Copy JavaScript files from TypeScript dist (excluding .d.ts files)
    echo "📦 Copying compiled JavaScript files..."
    cp {{project_root}}/ts/dist/*.js {{project_root}}/js/dist/ 2>/dev/null || true
    
    # Copy WASM files
    if [ -d "{{project_root}}/ts/dist/wasm" ]; then
        cp {{project_root}}/ts/dist/wasm/* {{project_root}}/js/dist/wasm/ 2>/dev/null || true
    elif [ -d "{{project_root}}/ts/wasm" ]; then
        cp {{project_root}}/ts/wasm/* {{project_root}}/js/dist/wasm/ 2>/dev/null || true
    fi
    
    echo "✅ JavaScript build complete!"

# Update JavaScript package: sync version, copy README and LICENSE
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
    
    # Ensure TypeScript package.json exists (use absolute path to be safe)
    if [ ! -f "{{project_root}}/ts/package.json" ]; then
        echo "❌ TypeScript package.json not found - cannot update JavaScript package"
        exit 1
    fi
    
    # Remove 'v' prefix if present for version
    VERSION="${CURRENT_TAG#v}"
    
    echo "📝 Generating JavaScript package.json from TypeScript..."
    cd {{project_root}}
    
    # Generate/update JavaScript package.json from TypeScript with git tag version
    node -e "
        const fs = require('fs');
        const tsPkg = JSON.parse(fs.readFileSync('ts/package.json', 'utf8'));
        
        // Create JavaScript package.json
        const jsPkg = {
            name: '@hlfshell/structured-parse-js',
            version: '$VERSION',
            description: tsPkg.description.replace('Multi-language', 'JavaScript').replace('TypeScript/JavaScript', 'JavaScript'),
            main: 'dist/index.js',
            type: 'module',
            scripts: {
                build: 'echo \"JavaScript build complete - files copied from TypeScript\"',
                prepublishOnly: 'npm run build'
            },
            keywords: tsPkg.keywords,
            author: tsPkg.author,
            license: tsPkg.license,
            files: [
                'dist/**/*',
                'README.md',
                'LICENSE'
            ]
        };
        
        // Ensure js directory exists
        if (!fs.existsSync('js')) {
            fs.mkdirSync('js', { recursive: true });
        }
        
        fs.writeFileSync('js/package.json', JSON.stringify(jsPkg, null, 2) + '\n');
    "
    
    echo "✅ JavaScript package.json generated!"
    
    # Copy README if it doesn't exist or is different
    if [ ! -f "{{project_root}}/js/README.md" ] || ! cmp -s "{{project_root}}/ts_js.md" "{{project_root}}/js/README.md" 2>/dev/null; then
        cp "{{project_root}}/ts_js.md" "{{project_root}}/js/README.md"
    fi
    
    # Copy LICENSE if it doesn't exist
    if [ ! -f "{{project_root}}/js/LICENSE" ]; then
        cp "{{project_root}}/LICENSE" "{{project_root}}/js/LICENSE"
    fi
    
    exit 0

# Publish JavaScript package to npm
publish: build
    #!/usr/bin/env bash
    set -euo pipefail
    
    # Ensure TypeScript package exists
    if [ ! -f "{{project_root}}/ts/package.json" ]; then
        echo "⚠️  TypeScript package.json not found - cannot publish JavaScript package"
        exit 0
    fi
        echo "📦 Publishing JavaScript package..."
    
    # Check if we're on a git tag
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before publishing."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    echo "✅ Current git tag: $CURRENT_TAG"
    
    # Ensure package.json exists and version matches git tag
    if [ ! -f "{{project_root}}/js/package.json" ]; then
        echo "❌ JavaScript package.json not found after build"
        exit 1
    fi
    
    # Verify version matches git tag (update should have set this, but double-check)
    VERSION="${CURRENT_TAG#v}"
    CURRENT_VERSION=$(cd {{project_root}}/js && node -p "require('./package.json').version" 2>/dev/null || echo "")
    if [ "$CURRENT_VERSION" != "$VERSION" ]; then
        echo "⚠️  Version mismatch: package.json has $CURRENT_VERSION but git tag is $VERSION"
        echo "   Updating version..."
        cd {{project_root}}/js
        node -e "const fs = require('fs'); const pkg = JSON.parse(fs.readFileSync('package.json', 'utf8')); pkg.version = '$VERSION'; fs.writeFileSync('package.json', JSON.stringify(pkg, null, 2) + '\n');"
    fi
    
    cd {{project_root}}/js
    npm publish --access public
    echo "✅ JavaScript publish complete!"

# Show available JavaScript commands
help:
    @echo "Available JavaScript commands:"
    @echo "  js clean    - Remove build artifacts and copied files"
    @echo "  js build    - Build JavaScript package (requires TypeScript build first)"
    @echo "  js test     - Run JavaScript tests"
    @echo "  js update   - Update version, README, and LICENSE, etc"
    @echo "  js publish  - Publish to npm"

