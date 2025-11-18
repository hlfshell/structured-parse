# structured-parse

**Parse labeled LLM output into clean, typed data structures.**  
Core implementation in Go, exported to TypeScript/JavaScript and Python via WebAssembly.

---

## Why structured-parse?

LLMs are good at *describing* data, but not always at emitting perfect JSON. `structured-parse` gives you a simple, robust format that LLMs can follow and your code can parse reliably.

If you have a prompt style like:

```text
Thought: ...
Action: ...
Action Input: {"query": "...", "max_results": 5}
Final Answer: ...
````

`structured-parse` turns it into predictable, typed data in your language of choice.

- **Flexible labels**
  - Case-insensitive matching
  - Multi-word labels
  - Multiple separators (`:`, `~`, `-`, `=` by default)
- **Validation**
  - Required fields
  - Cross-field dependencies (`requiredWith`)
  - Non-fatal, structured error reporting
- **Block parsing**
  - Parse multiple records from a single LLM response
- **JSON fields**
  - Mark specific labels as JSON and get parsed objects
- **Markdown aware**
  - Strips fenced code blocks and inline code markers before parsing
- **Multi-language**
  - Single Go implementation, shared with TypeScript/JavaScript and Python via WebAssembly

---

Each language has its own documentation, though they share the same API, differeing only in idiomaticy towards language styling. The library is written in GoLang, and then Python, Typescript and Javascript utilize WASM compilations to utilize the golang library.

You can learn more about your target version for each language at:

* golang - [go.md](go.md)
* Python - [python.md](python.md)
* Typescript / Javascript - [ts_js.md](ts_js.md)

---

## Installation

### Go

```bash
go get github.com/hlfshell/structured-parse/go/structuredparse
```

```go
import "github.com/hlfshell/structured-parse/go/structuredparse"
```

More in [go.md](go.md).

---

### TypeScript / JavaScript (Node.js)

```bash
npm install @hlfshell/structured-parse
# or
yarn add @hlfshell/structured-parse
```

```ts
import { createParser, type Label } from "@hlfshell/structured-parse";
```

More in [ts_js.md](ts_js.md).

---

### Python

`structured-parse` for Python uses the same WebAssembly core via a WASI runtime.

**Requirements:**

* Python 3.8+
* [`wasmtime` CLI](https://wasmtime.dev/) installed and on your `PATH`

Install:

```bash
pip install structured-parse
```

```python
from structured_parse import StructuredParser, Label
```

More in [python.md](python.md).

---

## Core concepts

### Labels

Labels describe the fields you expect from the LLM:

```go
type Label struct {
    Name         string   // Field name (case-insensitive match)
    Required     bool     // Must be present and non-empty
    RequiredWith []string // If this label is present, these must be too
    IsJSON       bool     // Parse the value as JSON
    IsBlockStart bool     // Marks the start of a new block/record
}
```

Key points:

* Matching is **case-insensitive** (`Action`, `action`, `ACTION` all match).
* Result keys use the **original `Name`** you provide (`"Action"`, `"Action Input"`, etc.).
* Unknown labels in the text are **ignored as separate fields** and treated as plain text.

### Separators

By default, `structured-parse` accepts the following separators between labels and values:

```text
Label: value
Label ~ value
Label - value
Label = value
```

You can restrict or customize separators through language-specific options (e.g. `ParserOptions.Separators` in Go).

### Multiline values

Values automatically span multiple lines until the next recognized label:

```text
Description: This is a long description
that spans multiple lines and will be
captured as a single value.
```

The parser will join these lines (with newlines preserved).

### Block parsing

If you set `IsBlockStart` on a label (for example, `Step`, `Task`, or `Name`), the parser will treat each occurrence as the start of a new “block” or record:

```text
Task: Task 1
Status: Complete
Result: Success

Task: Task 2
Status: In Progress
Result: Pending
```

You’ll get a list/array of maps/dicts/objects, one per block.

### JSON fields

Mark a label as JSON and get structured data:

```text
Config: {"retries": 3, "timeout": 10}
```

If JSON parsing fails:

* The **raw string** is kept as the value.
* An error is added to the error list.
* Parsing continues for other fields.

### Error handling

Errors are **non-fatal** by design:

* Missing required fields
* Failed `RequiredWith` dependencies
* JSON parse errors

Each language returns:

* A **result** object (map/dict) for successfully parsed data.
* A **list/array of error messages** you can log or handle.

---

## Performance

The parser is designed for typical LLM outputs (a few KB to tens of KB per response).

* Go benchmarks (`go/bench_test.go`) show sub-millisecond parsing for common patterns.
* TypeScript/JavaScript and Python share the same core engine via WASM and are suitable for request-per-response workloads.

If performance is critical in your environment, run the provided benchmarks and profile with your own prompts and label sets.

---

## License

MIT License – see [LICENSE](LICENSE) for details.
