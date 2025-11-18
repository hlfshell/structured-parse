"""Tests for the structured-parse parser."""

import json
import os
from pathlib import Path

try:
    import pytest
except ImportError:
    pytest = None

from structured_parse import StructuredParser, Label

# Get the path to test assets
ASSETS_DIR = Path(__file__).parent.parent.parent / "test-assets"


def deep_equal(a, b):
    """Deep equality check for dicts, lists, and primitives."""
    if type(a) != type(b):
        return False
    if isinstance(a, dict):
        if set(a.keys()) != set(b.keys()):
            return False
        return all(deep_equal(a[k], b[k]) for k in a.keys())
    if isinstance(a, list):
        if len(a) != len(b):
            return False
        return all(deep_equal(ai, bi) for ai, bi in zip(a, b))
    return a == b


def test_basic_functionality():
    """Test basic parsing functionality."""
    print("Testing basic functionality...")
    
    with open(ASSETS_DIR / "basic_functionality_input.txt") as f:
        input_text = f.read()
    with open(ASSETS_DIR / "basic_functionality_output.json") as f:
        expected = json.load(f)
    
    labels = [
        Label(name="Action Input", required_with=["Action"], is_json=True),
        Label(name="Action", required_with=["Action Input"]),
        Label(name="Thought"),
        Label(name="Result", required=True),
    ]
    
    parser = StructuredParser()
    result = parser.parse(labels, input_text)
    
    assert len(result.errors) == 0, f"Expected no errors, got: {result.errors}"
    assert deep_equal(result.result, expected), (
        f"Result mismatch:\nGot: {json.dumps(result.result, indent=2)}\n"
        f"Expected: {json.dumps(expected, indent=2)}"
    )
    
    print("✅ Basic functionality test passed")


def test_mixed_case_multiline():
    """Test mixed case and multiline values."""
    print("Testing mixed case and multiline...")
    
    with open(ASSETS_DIR / "mixed_case_multiline_input.txt") as f:
        input_text = f.read()
    with open(ASSETS_DIR / "mixed_case_multiline_output.json") as f:
        expected = json.load(f)
    
    labels = [
        Label(name="Context"),
        Label(name="Intention"),
        Label(name="Role"),
        Label(name="Action"),
        Label(name="Outcome"),
        Label(name="Notes"),
    ]
    
    parser = StructuredParser()
    result = parser.parse(labels, input_text)
    
    assert len(result.errors) == 0, f"Expected no errors, got: {result.errors}"
    assert deep_equal(result.result, expected), (
        f"Result mismatch:\nGot: {json.dumps(result.result, indent=2)}\n"
        f"Expected: {json.dumps(expected, indent=2)}"
    )
    
    print("✅ Mixed case multiline test passed")


def test_json_and_malformed():
    """Test JSON parsing and malformed JSON handling."""
    print("Testing JSON and malformed JSON...")
    
    with open(ASSETS_DIR / "json_and_malformed_input.txt") as f:
        input_text = f.read()
    with open(ASSETS_DIR / "json_and_malformed_output.json") as f:
        expected = json.load(f)
    with open(ASSETS_DIR / "json_and_malformed_errors.json") as f:
        expected_errors = json.load(f)
    
    labels = [
        Label(name="Config", is_json=True),
        Label(name="Data", is_json=True),
        Label(name="Description"),
    ]
    
    parser = StructuredParser()
    result = parser.parse(labels, input_text)
    
    assert deep_equal(result.result, expected), (
        f"Result mismatch:\nGot: {json.dumps(result.result, indent=2)}\n"
        f"Expected: {json.dumps(expected, indent=2)}"
    )
    assert len(result.errors) == len(expected_errors), (
        f"Expected {len(expected_errors)} errors, got {len(result.errors)}"
    )
    
    print("✅ JSON and malformed test passed")


def test_block_parsing():
    """Test block parsing."""
    print("Testing block parsing...")
    
    with open(ASSETS_DIR / "block_parsing_input.txt") as f:
        input_text = f.read()
    with open(ASSETS_DIR / "block_parsing_output.json") as f:
        expected = json.load(f)
    
    labels = [
        Label(name="Task", is_block_start=True),
        Label(name="Input", is_json=True),
        Label(name="Result"),
    ]
    
    parser = StructuredParser()
    result = parser.parse_blocks(labels, input_text)
    
    assert len(result.errors) == 0, f"Expected no errors, got: {result.errors}"
    assert deep_equal(result.blocks, expected), (
        f"Blocks mismatch:\nGot: {json.dumps(result.blocks, indent=2)}\n"
        f"Expected: {json.dumps(expected, indent=2)}"
    )
    
    print("✅ Block parsing test passed")


if __name__ == "__main__":
    # Run tests manually
    print("🧪 Running Python tests...\n")
    
    test_basic_functionality()
    test_mixed_case_multiline()
    test_json_and_malformed()
    test_block_parsing()
    
    print("\n✅ All Python tests passed!")


