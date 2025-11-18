# Go usage (`structuredparse`)

This document covers Go-specific usage of `structured-parse`.

Package import:

```go
import "github.com/hlfshell/structured-parse/go/structuredparse"
````

---

## Types

```go
type Label struct {
    Name         string
    Required     bool
    RequiredWith []string
    IsJSON       bool
    IsBlockStart bool
}

type ParserOptions struct {
    Separators string // Optional set of allowed separators (default: ":~-=")
}

type Parser struct {
    // constructed via NewParser
}
```

---

## Constructing a parser

```go
labels := []structuredparse.Label{
    {Name: "Thought"},
    {Name: "Action", RequiredWith: []string{"Action Input"}},
    {Name: "Action Input", IsJSON: true, RequiredWith: []string{"Action"}},
    {Name: "Final Answer"},
}

opts := &structuredparse.ParserOptions{
    // Use default separators (":~-=") by passing nil or zero-value options.
    // Separators: ":~-=",
}

parser, err := structuredparse.NewParser(labels, opts)
if err != nil {
    // e.g. duplicate labels, multiple block-start labels, invalid config
    log.Fatal(err)
}
```

Notes:

* `NewParser` **does not mutate** the `labels` slice you pass in.
* Label matching is case-insensitive, but result keys use the original `Name` values.

---

## Parsing a single record

```go
llmOutput := `
Reason: I see mostly positive language.
Sentiment: Positive
Confidence: {"score": 0.94, "threshold": 0.8}
`

labels := []structuredparse.Label{
    {Name: "Reason", Required: true},
    {Name: "Sentiment", Required: true},
    {Name: "Confidence", IsJSON: true},
}

parser, err := structuredparse.NewParser(labels, nil)
if err != nil {
    log.Fatal(err)
}

result, errs := parser.Parse(llmOutput)

if len(errs) > 0 {
    for _, e := range errs {
        log.Printf("parse warning: %s", e)
    }
}

// Typed access
reason, _ := result["Reason"].(string)
sentiment, _ := result["Sentiment"].(string)
confidence, _ := result["Confidence"].(map[string]interface{})

fmt.Println("Reason:", reason)
fmt.Println("Sentiment:", sentiment)
fmt.Println("Confidence:", confidence)
```

---

## Parsing multiple blocks

```go
llmOutput := `
Task: Data Collection
Status: Complete
Result: Success

Task: Trend Analysis
Status: In Progress
Result: Pending
`

labels := []structuredparse.Label{
    {Name: "Task", IsBlockStart: true, Required: true},
    {Name: "Status"},
    {Name: "Result"},
}

parser, err := structuredparse.NewParser(labels, nil)
if err != nil {
    log.Fatal(err)
}

blocks, errs := parser.ParseBlocks(llmOutput)

if len(errs) > 0 {
    for _, e := range errs {
        log.Printf("parse warning: %s", e)
    }
}

for i, b := range blocks {
    fmt.Printf("Block %d:\n", i+1)
    fmt.Printf("  Task:   %s\n", b["Task"])
    fmt.Printf("  Status: %s\n", b["Status"])
    fmt.Printf("  Result: %s\n\n", b["Result"])
}
```

Each `block` is a `map[string]interface{}` with keys matching your original label names.

---

## Custom separators

By default, these separators are accepted: `:`, `~`, `-`, `=`.

You can restrict or change them via `ParserOptions.Separators`:

```go
labels := []structuredparse.Label{
    {Name: "Key"},
    {Name: "Value"},
}

opts := &structuredparse.ParserOptions{
    Separators: ":", // only allow colon
}

parser, err := structuredparse.NewParser(labels, opts)
if err != nil {
    log.Fatal(err)
}
```

If a line doesn’t use a configured separator (e.g. `Key ~ value` when only `:` is allowed), it will not be recognized as a label line and will be treated as part of the current value instead.

---

## Multiline fields

Values automatically span multiple lines until the next recognized label:

```go
llmOutput := `
Description: This is a long description
that spans multiple lines and will be
captured as a single value.
Next Field: Done
`
```

`result["Description"]` will contain the entire multiline string (with newlines).

---

## Error handling

`Parse` and `ParseBlocks` return a result plus a `[]string` of errors:

* Missing required fields
* Failed `RequiredWith` dependencies
* JSON parse errors

Example:

```go
result, errs := parser.Parse(text)
if len(errs) > 0 {
    for _, e := range errs {
        log.Printf("parse warning: %s", e)
    }
}

// Even if there are warnings, result may still be partially or wholly usable.
```

---

## Markdown handling

Before parsing, `structured-parse` strips:

* Fenced code blocks:

  ````text
  ```json
  {"key":"value"}
  ````

  ```
  ```
* Inline code:

  ```text
  `some inline code`
  ```

This helps when your prompts wrap JSON or examples in markdown for readability.

If you need different behavior, you can pre-process the text before passing it to `Parse` / `ParseBlocks`.
