
## `ts_js.md`

```markdown
# TypeScript / JavaScript usage

The TypeScript/JavaScript bindings call into the core Go implementation compiled to WebAssembly. Initialization is asynchronous.

## Package Selection

We provide two npm packages:

### TypeScript Package (`@hlfshell/structured-parse`)

```bash
npm install @hlfshell/structured-parse
# or
yarn add @hlfshell/structured-parse
pnpm add @hlfshell/structured-parse
```

This package includes:
- Full TypeScript type definitions
- TypeScript source code
- Compiled JavaScript output

### JavaScript Package (`@hlfshell/structured-parse-js`)

```bash
npm install @hlfshell/structured-parse-js
# or
yarn add @hlfshell/structured-parse-js
pnpm add @hlfshell/structured-parse-js
```

---

## API overview

### TypeScript

```ts
import { createParser, type Label, type ParseResult, type ParseBlocksResult } from "@hlfshell/structured-parse";
```

### JavaScript

```js
import { createParser } from "@hlfshell/structured-parse-js";
```

interface Label {
  name: string;
  required?: boolean;
  requiredWith?: string[];
  isJson?: boolean;
  isBlockStart?: boolean;
}

interface ParserOptions {
  separators?: string; // default: ":~-="
}

interface ParseResult {
  result: Record<string, unknown>;
  errors: string[];
}

interface ParseBlocksResult {
  blocks: Array<Record<string, unknown>>;
  errors: string[];
}

interface Parser {
  parse(labels: Label[], text: string, options?: ParserOptions): ParseResult;
  parseBlocks(labels: Label[], text: string, options?: ParserOptions): ParseBlocksResult;
}
```

---

## Initialization

Because the core is WebAssembly, you must `await` parser creation:

### TypeScript

```ts
import { createParser } from "@hlfshell/structured-parse";
```

### JavaScript

```js
import { createParser } from "@hlfshell/structured-parse-js";
```

async function main() {
  const parser = await createParser();
  // use parser...
}

main().catch(console.error);
```

---

## Parsing a single record

### TypeScript

```ts
import { createParser, type Label } from "@hlfshell/structured-parse";
```

async function parseSentiment(llmOutput: string) {
  const parser = await createParser();

  const labels: Label[] = [
    { name: "Reason", required: true },
    { name: "Sentiment", required: true },
    { name: "Confidence", isJson: true },
  ];

  const { result, errors } = parser.parse(labels, llmOutput);

  if (errors.length > 0) {
    console.warn("parse warnings:", errors);
  }

  const reason = result["Reason"] as string | undefined;
  const sentiment = result["Sentiment"] as string | undefined;
  const confidence = result["Confidence"] as { score: number; threshold: number } | undefined;

  console.log("Reason:", reason);
  console.log("Sentiment:", sentiment);
  console.log("Confidence:", confidence);
}
```

### JavaScript

```js
import { createParser } from "@hlfshell/structured-parse-js";

async function parseSentiment(llmOutput) {
  const parser = await createParser();

  const labels = [
    { name: "Reason", required: true },
    { name: "Sentiment", required: true },
    { name: "Confidence", isJson: true },
  ];

  const { result, errors } = parser.parse(labels, llmOutput);

  if (errors.length > 0) {
    console.warn("parse warnings:", errors);
  }

  const reason = result["Reason"];
  const sentiment = result["Sentiment"];
  const confidence = result["Confidence"];

  console.log("Reason:", reason);
  console.log("Sentiment:", sentiment);
  console.log("Confidence:", confidence);
}
```

Notes:

* Label matching is case-insensitive, but `result` keys use the original `name` values.
* JSON fields (`isJson: true`) are parsed via `JSON.parse` in the core; on failure, you get:

  * The raw string as the value
  * An error message in `errors`

---

## Parsing multiple blocks

### TypeScript

```ts
import { createParser, type Label } from "@hlfshell/structured-parse";
```

async function parseBlocksExample(llmOutput: string) {
  const parser = await createParser();

  const labels: Label[] = [
    { name: "Task", isBlockStart: true, required: true },
    { name: "Status" },
    { name: "Result" },
  ];

  const { blocks, errors } = parser.parseBlocks(labels, llmOutput);

  if (errors.length > 0) {
    console.warn("parse warnings:", errors);
  }

  blocks.forEach((block, index) => {
    console.log(`Block ${index + 1}:`);
    console.log("  Task:", block["Task"]);
    console.log("  Status:", block["Status"]);
    console.log("  Result:", block["Result"]);
  });
}
```

### JavaScript

```js
import { createParser } from "@hlfshell/structured-parse-js";

async function parseBlocksExample(llmOutput) {
  const parser = await createParser();

  const labels = [
    { name: "Task", isBlockStart: true, required: true },
    { name: "Status" },
    { name: "Result" },
  ];

  const { blocks, errors } = parser.parseBlocks(labels, llmOutput);

  if (errors.length > 0) {
    console.warn("parse warnings:", errors);
  }

  blocks.forEach((block, index) => {
    console.log(`Block ${index + 1}:`);
    console.log("  Task:", block["Task"]);
    console.log("  Status:", block["Status"]);
    console.log("  Result:", block["Result"]);
  });
}
```

Each `block` is a `Record<string, unknown>` (TypeScript) or plain object (JavaScript).

---

## Custom separators

Default separators: `:`, `~`, `-`, `=`.

You can override them:

### TypeScript

```ts
const labels: Label[] = [
  { name: "Key" },
  { name: "Value" },
];

const options = { separators: ":" }; // only colon

const { result, errors } = parser.parse(labels, "Key: foo\nValue: bar", options);
```

### JavaScript

```js
const labels = [
  { name: "Key" },
  { name: "Value" },
];

const options = { separators: ":" }; // only colon

const { result, errors } = parser.parse(labels, "Key: foo\nValue: bar", options);
```

If a line doesn't use a configured separator, it will not be recognized as a label line and will be treated as part of the current field's value.

---

## Multiline values

Values span multiple lines until the next recognized label:

### TypeScript

```ts
const llmOutput = `
Description: This is a long description
that spans multiple lines and will be
captured as a single value.
Next Field: Done
`;

const labels: Label[] = [
  { name: "Description" },
  { name: "Next Field" },
];

const { result } = parser.parse(labels, llmOutput);
console.log(result["Description"]);
// "This is a long description\nthat spans multiple lines and will be\ncaptured as a single value."
```

### JavaScript

```js
const llmOutput = `
Description: This is a long description
that spans multiple lines and will be
captured as a single value.
Next Field: Done
`;

const labels = [
  { name: "Description" },
  { name: "Next Field" },
];

const { result } = parser.parse(labels, llmOutput);
console.log(result["Description"]);
// "This is a long description\nthat spans multiple lines and will be\ncaptured as a single value."
```

---

## Error handling

Both `parse` and `parseBlocks` return an `errors: string[]` array:

```ts
const { result, errors } = parser.parse(labels, text);

if (errors.length > 0) {
  console.warn("parse warnings:", errors);
}

// result may still be fully or partially usable
```

Errors include:

* Missing required fields (e.g. `"Sentiment" is required`)
* Failed `requiredWith` dependencies (e.g. `"Action" requires "Action Input"`)
* JSON parse failures (e.g. `JSON error in 'Config': ...`)

---

## Browser usage (high-level)

The package can be used in browser environments that support WebAssembly. Typical setup:

* Bundle the WASM and JS with your tool (Vite, Webpack, etc.).
* Initialize the parser via `await createParser()` in your app code.

Example (simplified):

### TypeScript

```ts
import { createParser } from "@hlfshell/structured-parse";
```

### JavaScript

```js
import { createParser } from "@hlfshell/structured-parse-js";
```

let parserPromise: Promise<ReturnType<typeof createParser>> | null = null;

export function getParser() {
  if (!parserPromise) {
    parserPromise = createParser();
  }
  return parserPromise;
}

// elsewhere
const parser = await getParser();
const { result } = parser.parse(labels, text);
```

Consult your bundler’s documentation for handling `.wasm` assets if needed.
