# structured-parse Specification

This document provides detailed technical specifications for the `structured-parse` library. It defines the exact behavior, parsing rules, and API contracts that all implementations must follow.

## Overview

`structured-parse` is a library for parsing labeled/structured LLM output into structured data. It supports:
- Case-insensitive label matching with case-preserving result keys
- Multi-line value collection
- JSON field parsing
- Required field and dependency validation
- Block-based parsing for multiple entries

## Label Definition

A `Label` is defined with the following properties:

- **Name** (string): The label name to match. Matching is case-insensitive, but the original casing is preserved in result map keys. Multi-word labels are supported (e.g., "Action Input").
- **Required** (bool): If `true`, this label must be present in the input. Missing required labels generate an error.
- **RequiredWith** ([]string): List of other label names that must also be present if this label is present. Dependency is enforced only if the label itself is present (even if empty).
- **IsJSON** (bool): If `true`, the label value is parsed as JSON. Empty JSON fields become `{}`. Invalid JSON is kept as a string and an error is added.
- **IsBlockStart** (bool): If `true`, this label marks the start of a new block for block parsing. Exactly one label must have `IsBlockStart: true`.

### Label Name Constraints

- Label names are matched case-insensitively
- Multi-word labels are supported (whitespace is normalized in matching)
- Original casing provided by the user is preserved in result map keys
- Label names in `RequiredWith` are also matched case-insensitively

## Input Preprocessing

Before parsing, the input text undergoes the following cleaning steps:

1. **Markdown Code Block Removal**: All markdown code blocks (```...```) are removed. The content inside the code block is extracted and used as the text content.
   - Pattern: `(?s)```(?:\\w+)?\\s*(.*?)\\s*````
   - The language identifier (if present) is ignored
   - Only the content between the triple backticks is kept

2. **Inline Code Removal**: All inline code markers (`...`) are removed, keeping only the content.
   - Pattern: `` `([^`]+)` ``
   - The backticks are removed, content is preserved

3. **Whitespace Trimming**: Leading and trailing whitespace is removed from the entire text.

4. **Line Processing**: The text is split into lines. Each line has its right-side whitespace (spaces, tabs, carriage returns) trimmed.

## Label Matching Rules

### Matching Algorithm

Labels are matched using a two-phase approach:

1. **Primary Matching (Regex)**: For each label, a regex pattern is constructed:
   - Multi-word labels have whitespace normalized to `\s+`
   - Pattern: `(?i)^\s*{label}\s*[:~\-]+\s*`
   - The pattern is case-insensitive (`(?i)`)
   - Allows optional leading whitespace
   - Matches one or more separators: colon (`:`), tilde (`~`), or dash (`-`)
   - Allows optional trailing whitespace before the value

2. **Fallback Matching**: If regex matching fails, a prefix-based check is performed:
   - Line is trimmed of leading whitespace
   - Case-insensitive prefix match against label name
   - Remaining text must start with a separator pattern: `^\s*[:~\-]+`
   - If no separator is found, the line is treated as continuation text

### Matching Behavior

- Labels must appear at the **start of a line** (after optional leading whitespace)
- Matching is **case-insensitive**: `Task:`, `task:`, `TASK:` all match the same label
- Multiple separators are supported: `:`, `~`, `-` (with optional whitespace)
- The separator and any trailing whitespace are stripped from the value
- If a line does not match any label, it is treated as continuation text for the current label

### Unknown Labels

Lines that look like labels (have a label-like pattern) but don't match any defined label are **ignored as distinct fields**. They are treated as continuation text for the current field being collected.

## Value Collection

### Single-Line Values

If a label is followed by a value on the same line:
```
Label: value
```
The value is the text after the separator (trimmed of whitespace).

### Multi-Line Values

Values can span multiple lines. Continuation rules:

1. A line continues the current label's value if:
   - It does not start with a known label (case-insensitive check)
   - The previous line was collecting a value

2. Continuation lines are joined with newline characters (`\n`)

3. If a continuation line starts with a known label name (even if not properly formatted), it is treated as a new label, not continuation.

### Empty Values

- If a label appears with no value: `Label:`, the value is an empty string `""`
- If a label is defined but never appears in the input, its value in the result is `""`
- Empty values are preserved (not filtered out)

### Multiple Occurrences

If a label appears multiple times in the input, all values are collected into a slice. The result map will contain a slice of values for that label.

## Result Map Structure

### Key Naming

- **Result map keys use the original `Label.Name` casing** as provided by the user
- Internal matching is case-insensitive, but keys preserve original casing
- Example: If label is defined as `{Name: "Reason"}`, the result key is `"Reason"`, not `"reason"`

### Value Types

Each value in the result map can be:

1. **String**: For plain text values (single occurrence, non-JSON)
2. **Empty String** (`""`): For missing or empty labels
3. **Parsed JSON Object**: For labels marked with `IsJSON: true` (when JSON parsing succeeds)
4. **String** (fallback): For labels marked with `IsJSON: true` but with invalid JSON (raw string is kept)
5. **Slice**: If a label appears multiple times, values are collected into a slice `[]interface{}`

### Value Flattening

- **Single value**: If a label appears once, the value is stored directly (not in a slice)
- **Multiple values**: If a label appears multiple times, values are stored in a slice
- **Empty values**: Empty strings are flattened to `""`, not `[""]`

## JSON Field Parsing

### JSON Parsing Rules

For labels marked with `IsJSON: true`:

1. **Empty JSON Fields**: If the value is empty or only whitespace, it becomes an empty object `{}`
2. **Valid JSON**: The value is parsed using standard JSON parsing. The result can be any valid JSON type (object, array, string, number, boolean, null)
3. **Invalid JSON**: If parsing fails:
   - The raw string value is kept in the result
   - An error is added to the error list: `"JSON error in '{LabelName}': {error message}"`
   - The error message uses the original label name (preserving casing)

### JSON Error Handling

- JSON parsing errors do not stop parsing
- The raw string is preserved for debugging
- Errors are collected and returned separately

## Required Field Validation

### Required Labels

If a label has `Required: true` and is missing from the input, an error is generated:
- Format: `"{LabelName}' is required"`
- Uses original label name casing

### Missing Field Detection

A field is considered "missing" if:
- The label never appears in the input, OR
- The label appears but has an empty value (empty string after trimming)

### RequiredWith Dependencies

If a label has `RequiredWith: ["OtherLabel"]`:

- The dependency is **only enforced if the label itself is present** (even if empty)
- If the label is present but the dependency is missing, an error is generated:
  - Format: `"{LabelName}' requires '{DependencyName}'"`
  - Uses original label names (preserving casing)
- Dependency names in `RequiredWith` are matched case-insensitively

## Block Parsing

### Block Definition

Blocks are used to parse multiple independent entries from a single input. Blocks are separated by the label marked with `IsBlockStart: true`.

### Block Parsing Rules

1. **Block Start Label**: Exactly one label must have `IsBlockStart: true`. If none is defined, `ParseBlocks` returns an error.

2. **Block Splitting**: The input is split into blocks at each occurrence of the block start label:
   - Each block starts with the block start label
   - Blocks are independent (each is parsed separately)
   - Lines before the first block start label are ignored

3. **Block Parsing**: Each block is parsed using the standard `Parse` logic:
   - Each block is treated as a separate document
   - All labels (including the block start label) can appear in each block
   - Validation (required fields, dependencies) is performed per block

4. **Error Collection**: Errors from all blocks are collected into a single error list

5. **Result Structure**: Returns a slice of result maps, one per block

### Block Start Label Matching

The block start label is matched using the same case-insensitive rules as regular labels. The label name used for matching is the lowercase version, but the original casing is preserved in block result maps.

## Error Handling

### Error Format

Errors are returned as a slice of strings (`[]string`). Each error is a human-readable message.

### Error Types

1. **Missing Required Field**: `"{LabelName}' is required"`
2. **Missing Dependency**: `"{LabelName}' requires '{DependencyName}'"`
3. **JSON Parsing Error**: `"JSON error in '{LabelName}': {JSON error message}"`
4. **Block Parsing Error**: `"no block start label defined - must have at least one"` (returned by `ParseBlocks` if no block start label is defined)

### Error Message Casing

All error messages use the **original label name casing** as provided by the user, not the internal lowercase version.

## API Contracts

### NewParser

```go
NewParser(labels []Label, opts *ParserOptions) (*Parser, error)
```

**Behavior:**
- Creates a new parser with the given labels and optional options
- If `opts` is `nil`, default options are used (separators: `":~-="`)
- If `opts` is provided, custom separators can be specified via `opts.Separators`
- **Does not modify the input `labels` slice** (makes an internal copy)
- Returns an error if more than one label has `IsBlockStart: true`
- Error message: `"only one block start label is allowed"` (lowercase, no trailing period)

**ParserOptions:**
- `Separators` (string): Allowed separator characters. Default is `":~-="` (colon, tilde, dash, equals). Each character in the string is treated as a valid separator.

**Thread Safety:**
- Parser construction is safe for concurrent use
- The returned `Parser` is safe for concurrent use (read-only configuration)

### Parse

```go
Parse(text string) (map[string]interface{}, []string)
```

**Behavior:**
- Parses a single document from the input text
- Returns a result map and a slice of error strings
- Result map keys use original label name casing
- Empty labels result in `""` values
- Multiple occurrences result in slices

### ParseBlocks

```go
ParseBlocks(text string) ([]map[string]interface{}, []string)
```

**Behavior:**
- Parses multiple blocks from the input text
- Requires exactly one label with `IsBlockStart: true`
- Returns a slice of result maps (one per block) and a combined error list
- If no block start label is defined, returns `nil, []string{"no block start label defined - must have at least one"}`

## Edge Cases

### Empty Input

- Empty input results in all labels having `""` values
- Required field validation still runs (will generate errors if labels are required)

### Whitespace-Only Input

- Whitespace-only input is treated as empty
- All labels get `""` values

### Labels with No Separator

- If a label name appears at the start of a line but without a valid separator, it is treated as continuation text (not a new label)

### Overlapping Label Names

- If one label name is a prefix of another (e.g., "Action" and "Action Input"), the longer match takes precedence in regex matching
- Fallback matching may match the shorter label first

### Special Characters in Values

- Newlines, tabs, and other special characters are preserved in values
- JSON fields must contain valid JSON (special characters must be escaped)

### Markdown in Values

- Markdown code blocks and inline code in values are **not** removed (only removed from the overall input during preprocessing)
- If a value itself contains triple backticks, they are preserved

## Implementation Notes

### Internal Data Structures

- Labels are stored internally with lowercase names for matching
- A mapping from lowercase names to original names is maintained for result keys
- Regex patterns are compiled once during parser construction

### Performance Considerations

- Regex patterns are compiled once per parser instance
- Label matching uses compiled regex patterns for efficiency
- Multi-line value collection uses a string builder for efficiency
- JSON parsing is performed per occurrence (not batched)

### Concurrency

- `Parser` instances are safe for concurrent use
- Multiple goroutines can call `Parse` or `ParseBlocks` on the same `Parser` instance concurrently
- No shared mutable state exists between parse operations

