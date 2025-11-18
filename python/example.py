"""
Example usage of structured-parse Python library
"""

import json
from structured_parse import StructuredParser, Label


def main():
    print("🚀 structured-parse Python Example\n")

    # Create the parser
    print("Initializing parser (loading WASM)...")
    parser = StructuredParser()
    print("✅ Parser ready!\n")

    # Example 1: Basic parsing with JSON field
    print("Example 1: Basic Parse")
    print("─" * 50)

    labels1 = [
        Label(name="Action", required=True),
        Label(name="Parameters", is_json=True),
        Label(name="Thought"),
    ]

    text1 = """
Action: process_data
Parameters: {"input": ["file1.txt", "file2.txt"], "mode": "batch"}
Thought: Processing multiple files in batch mode
    """.strip()

    result1 = parser.parse(labels1, text1)

    print(f"Input: {text1}")
    print(f"\nParsed Result:")
    print(json.dumps(result1.result, indent=2))
    print(f"Errors: {'None' if not result1.errors else result1.errors}")

    # Example 2: Block parsing
    print("\n\nExample 2: Block Parsing")
    print("─" * 50)

    labels2 = [
        Label(name="Task", is_block_start=True),
        Label(name="Input", is_json=True),
        Label(name="Result"),
    ]

    text2 = """
Task: Summarize
Input: {"text": "Long article about AI..."}
Result: AI article discusses recent developments

Task: Classify
Input: {"text": "Product review..."}
Result: Positive sentiment
    """.strip()

    result2 = parser.parse_blocks(labels2, text2)

    print(f"Input: {text2}")
    print(f"\nParsed Blocks:")
    for i, block in enumerate(result2.blocks, 1):
        print(f"\nBlock {i}:")
        print(json.dumps(block, indent=2))
    print(f"Errors: {'None' if not result2.errors else result2.errors}")

    # Example 3: Validation (required field missing)
    print("\n\nExample 3: Validation Errors")
    print("─" * 50)

    labels3 = [
        Label(name="Name", required=True),
        Label(name="Email", required_with=["Name"]),
    ]

    text3 = "Email: test@example.com"  # Missing required 'Name'

    result3 = parser.parse(labels3, text3)

    print(f"Input: {text3}")
    print(f"\nParsed Result: {json.dumps(result3.result, indent=2)}")
    print(f"Errors: {result3.errors}")

    print("\n✅ Examples complete!")


if __name__ == "__main__":
    main()


