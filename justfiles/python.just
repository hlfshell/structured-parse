# Python-specific tasks for structured-parse

# Use invocation directory (where just was called from) as project root
project_root := invocation_directory()

# Default recipe - show help
default: help

# Install Python dependencies
install:
    @echo "📦 Installing Python dependencies..."
    @eval "$(pyenv init -)" && eval "$(pyenv virtualenv-init -)" && pyenv activate structured-parse && cd {{project_root}}/python && pip install -e .[dev]
    @echo "✅ Python dependencies installed!"

# Build Python package
build:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔨 Building Python package..."
    eval "$(pyenv init -)" && eval "$(pyenv virtualenv-init -)" && pyenv activate structured-parse
    if [ -f "{{project_root}}/python/pyproject.toml" ]; then
        cd {{project_root}}/python && python -m build
    elif [ -f "{{project_root}}/python/setup.py" ]; then
        cd {{project_root}}/python && python setup.py sdist bdist_wheel
    else
        echo "❌ Python build configuration not found (pyproject.toml or setup.py)"
        exit 1
    fi
    echo "✅ Python build complete!"

# Run Python tests
test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🧪 Running Python tests..."
    eval "$(pyenv init -)" && eval "$(pyenv virtualenv-init -)" && pyenv activate structured-parse
    if [ -f "{{project_root}}/python/pyproject.toml" ] || [ -f "{{project_root}}/python/setup.py" ]; then
        cd {{project_root}}/python
        if command -v pytest >/dev/null 2>&1; then
            python -m pytest tests/ -v
        elif python -c "import pytest" 2>/dev/null; then
            python -m pytest tests/ -v
        else
            echo "⚠️  pytest not installed. Running tests directly..."
            python tests/test_parser.py
        fi
    else
        echo "❌ Python test configuration not found"
        exit 1
    fi
    echo "✅ Python tests complete!"

# Publish Python package to PyPI
publish:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "📦 Publishing Python package..."
    eval "$(pyenv init -)" && eval "$(pyenv virtualenv-init -)" && pyenv activate structured-parse
    
    # Check if we're on a git tag
    CURRENT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "")
    if [ -z "$CURRENT_TAG" ]; then
        echo "❌ Error: Not on a git tag. Please checkout a tag before publishing."
        echo "   Example: git checkout v1.0.0"
        exit 1
    fi
    
    echo "✅ Current git tag: $CURRENT_TAG"
    
    if [ -f "{{project_root}}/python/pyproject.toml" ] || [ -f "{{project_root}}/python/setup.py" ]; then
        cd {{project_root}}/python
        if command -v twine &> /dev/null; then
            python -m build 2>/dev/null || python setup.py sdist bdist_wheel
            twine upload dist/*
        else
            echo "❌ twine not installed. Install with: pip install twine"
            exit 1
        fi
    else
        echo "❌ Python publish configuration not found"
        exit 1
    fi
    echo "✅ Python publish complete!"

# Show available Python commands
help:
    @echo "Available Python commands:"
    @echo "  python install  - Install dependencies"
    @echo "  python build    - Build Python package"
    @echo "  python test     - Run Python tests"
    @echo "  python publish  - Publish to PyPI (requires git tag)"

