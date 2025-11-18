# Go-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Build Go package
build:
    @echo "🔨 Building Go..."
    @cd {{project_root}}/go && go build ./...
    @echo "✅ Go build complete!"

# Run Go tests
test:
    @echo "🧪 Running Go tests..."
    @cd {{project_root}}/go && go test -v ./...
    @echo "✅ Go tests complete!"

# Run Go benchmarks
bench:
    @echo "📊 Running Go benchmarks..."
    @cd {{project_root}}/go && go test -bench=. -benchmem ./...
    @echo "✅ Go benchmarks complete!"

# Build WASM module for JS/TS (using standard Go WASM)
build-wasm-js:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔨 Building WASM module for JS/TS..."
    mkdir -p {{project_root}}/ts/wasm {{project_root}}/js/wasm
    cd {{project_root}}/go/cmd/wasm && GOOS=js GOARCH=wasm go build -o {{project_root}}/ts/wasm/structured-parse.wasm
    GOROOT=$(go env GOROOT)
    # Try multiple possible locations for wasm_exec.js
    if [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
        cp "$GOROOT/misc/wasm/wasm_exec.js" {{project_root}}/ts/wasm/
    elif [ -f "$GOROOT/lib/wasm/wasm_exec.js" ]; then
        cp "$GOROOT/lib/wasm/wasm_exec.js" {{project_root}}/ts/wasm/
    elif [ -f "/usr/share/go/misc/wasm/wasm_exec.js" ]; then
        cp "/usr/share/go/misc/wasm/wasm_exec.js" {{project_root}}/ts/wasm/
    else
        echo "❌ Could not find wasm_exec.js. Trying to find it..."
        WASM_EXEC=$(find $GOROOT -name "wasm_exec.js" 2>/dev/null | head -1)
        if [ -z "$WASM_EXEC" ]; then
            echo "❌ wasm_exec.js not found in Go installation"
            exit 1
        fi
        cp "$WASM_EXEC" {{project_root}}/ts/wasm/
    fi
    cp {{project_root}}/ts/wasm/structured-parse.wasm {{project_root}}/js/wasm/
    cp {{project_root}}/ts/wasm/wasm_exec.js {{project_root}}/js/wasm/
    echo "✅ JS/TS WASM build complete!"

# Build WASM module for Python (using TinyGo for WASI compatibility)
build-wasm-python:
    @echo "🔨 Building WASM module for Python..."
    @mkdir -p {{project_root}}/python/structured_parse/wasm
    @if command -v tinygo >/dev/null 2>&1; then \
        cd {{project_root}}/go/cmd/wasm-wasi && \
        tinygo build -o {{project_root}}/python/structured_parse/wasm/structured-parse.wasm -target=wasi .; \
        echo "✅ Python WASM build complete!"; \
    else \
        echo "⚠️  TinyGo not found. Installing TinyGo is required for Python WASM builds."; \
        echo "   Visit: https://tinygo.org/getting-started/install/"; \
        cd {{project_root}}/go/cmd/wasm-wasi && \
        GOOS=wasip1 GOARCH=wasm go build -o {{project_root}}/python/structured_parse/wasm/structured-parse.wasm .; \
        echo "✅ Python WASM build complete (using Go WASI)!"; \
    fi

# Build all WASM modules
build-wasm: build-wasm-js build-wasm-python
    @echo "✅ All WASM builds complete!"

# Show available Go commands
help:
    @echo "Available Go commands:"
    @echo "  go build           - Build Go package"
    @echo "  go test            - Run Go tests"
    @echo "  go bench           - Run Go benchmarks"
    @echo "  go build-wasm-js   - Build WASM for JS/TS"
    @echo "  go build-wasm-python - Build WASM for Python"
    @echo "  go build-wasm      - Build all WASM modules"

