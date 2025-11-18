# Python usage

The Python bindings call into the core Go implementation compiled to WebAssembly using a WASI runtime.

## Requirements

- Python 3.8+
- [`wasmtime` CLI](https://wasmtime.dev/) installed and available on your `PATH`

Install:

```bash
pip install structured-parse
````

Basic import:

```python
from structured_parse import StructuredParser, Label
```

---

## Types

The exact implementation may vary, but conceptually:

```python
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Label:
    name: str
    required: bool = False
    required_with: Optional[List[str]] = None
    is_json: bool = False
    is_block_start: bool = False

@dataclass
class ParserOptions:
    separators: str = ":~-="  # default separators

@dataclass
class ParseResult:
    result: dict[str, object]
    errors: list[str]

@dataclass
class ParseBlocksResult:
    blocks: list[dict[str, object]]
    errors: list[str]
```

`StructuredParser` uses the WASM module under the hood; the interface is synchronous, but initialization may perform some one-time setup on first use.

---

## Creating a parser

```python
from structured_parse import StructuredParser, Label

labels = [
    Label(name="Reason", required=True),
    Label(name="Sentiment", required=True),
    Label(name="Confidence", is_json=True),
]

parser = StructuredParser()  # default options (":~-=" separators)
```

If you expose options:

```python
from structured_parse import StructuredParser, Label, ParserOptions

options = ParserOptions(separators=":")  # only allow colon
parser = StructuredParser(options=options)
```

---

## Parsing a single record

```python
from structured_parse import StructuredParser, Label

labels = [
    Label(name="Reason", required=True),
    Label(name="Sentiment", required=True),
    Label(name="Confidence", is_json=True),
]

llm_output = """
Reason: I see mostly positive language.
Sentiment: Positive
Confidence: {"score": 0.94, "threshold": 0.8}
"""

parser = StructuredParser()
result = parser.parse(labels, llm_output)

if result.errors:
    print("Warnings:")
    for err in result.errors:
        print("  -", err)

reason = result.result.get("Reason")
sentiment = result.result.get("Sentiment")
confidence = result.result.get("Confidence")  # likely a dict

print("Reason:", reason)
print("Sentiment:", sentiment)
print("Confidence:", confidence)
```

Notes:

* Matching is case-insensitive, but `result.result` keys use your original label names (`"Reason"`, `"Sentiment"`, etc.).
* JSON fields (`is_json=True`) are parsed; on failure:

  * The raw string is kept.
  * An error is added to `result.errors`.

---

## Parsing multiple blocks

```python
from structured_parse import StructuredParser, Label

labels = [
    Label(name="Step", is_block_start=True, required=True),
    Label(name="Analysis"),
    Label(name="Data", is_json=True),
    Label(name="Conclusion"),
]

llm_output = """
Step: Data Collection
Analysis: Gathered user feedback from 1,000 responses
Data: {"positive": 650, "neutral": 250, "negative": 100}
Conclusion: Majority positive sentiment

Step: Trend Analysis
Analysis: Comparing with previous quarter
Data: {"growth": 15.5, "retention": 92.3}
Conclusion: Strong upward trend
"""

parser = StructuredParser()
blocks_result = parser.parse_blocks(labels, llm_output)

for i, block in enumerate(blocks_result.blocks, start=1):
    print(f"\n=== Step {i}: {block.get('Step')} ===")
    print("Analysis:", block.get("Analysis"))
    print("Data:", block.get("Data"))
    print("Conclusion:", block.get("Conclusion"))

if blocks_result.errors:
    print("\nWarnings:")
    for err in blocks_result.errors:
        print("  -", err)
```

Each block is a `dict[str, object]`.

---

## Custom separators

Default separators: `:`, `~`, `-`, `=`.

If `ParserOptions` is exposed, you can restrict or change them:

```python
from structured_parse import StructuredParser, Label, ParserOptions

labels = [
    Label(name="Key"),
    Label(name="Value"),
]

options = ParserOptions(separators=":")  # only colon
parser = StructuredParser(options=options)

output = """
Key: foo
Value: bar
"""

result = parser.parse(labels, output)
print(result.result)
```

Lines without a configured separator aren’t recognized as label lines and are treated as part of the current field’s value.

---

## Multiline values

Values span multiple lines until the next recognized label:

```python
labels = [
    Label(name="Description"),
    Label(name="Next Field"),
]

llm_output = """
Description: This is a long description
that spans multiple lines and will be
captured as a single value.
Next Field: Done
"""

parser = StructuredParser()
result = parser.parse(labels, llm_output)

print(result.result["Description"])
# "This is a long description\nthat spans multiple lines and will be\ncaptured as a single value."
```

---

## Error handling

Both `parse` and `parse_blocks` return result objects with an `errors` list:

```python
result = parser.parse(labels, text)

if result.errors:
    print("Warnings:")
    for err in result.errors:
        print("  -", err)

# result.result is available even if there are warnings
```

Common errors:

* Missing required fields: `"Sentiment" is required`
* Dependency failures: `"Action" requires "Action Input"`
* JSON parse errors: `JSON error in 'Config': ...`

---

## WASM and wasmtime

Under the hood, `StructuredParser`:

1. Loads the WebAssembly module built from the Go implementation.
2. Executes it via the `wasmtime` CLI (WASI).

This means:

* You get the same behavior and bugfixes as the Go implementation.
* You must have `wasmtime` installed and accessible on your system.
* If initialization fails (e.g., WASM module not found, wasmtime missing), the library will raise a Python exception.

If you encounter WASM-related issues, verify:

* `wasmtime --version` works in your shell.
* Your environment can find the WASM module packaged with `structured-parse`.
