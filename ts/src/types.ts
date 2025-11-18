/**
 * Label defines a label for parsing with options for required, dependencies, JSON, and block start.
 */
export interface Label {
  /** Name of the label (case-insensitive matching, original casing preserved in results) */
  name: string;
  /** Whether this label is required */
  required?: boolean;
  /** List of other label names required with this one */
  requiredWith?: string[];
  /** Whether this label should be parsed as JSON */
  isJson?: boolean;
  /** Whether this label starts a new block */
  isBlockStart?: boolean;
}

/**
 * ParserOptions allows customization of parser behavior.
 */
export interface ParserOptions {
  /** Allowed separator characters. Default is ":~-=" */
  separators?: string;
}

/**
 * ParseResult contains the parsed data and any errors encountered.
 */
export interface ParseResult {
  /** The parsed result map (label name -> value) */
  result: Record<string, any>;
  /** List of parsing/validation errors */
  errors: string[];
}

/**
 * ParseBlocksResult contains the parsed blocks and any errors encountered.
 */
export interface ParseBlocksResult {
  /** Array of parsed blocks, each block is a map of label name -> value */
  blocks: Array<Record<string, any>>;
  /** List of parsing/validation errors */
  errors: string[];
}

/**
 * Internal WASM response structure
 */
export interface WasmResponse {
  ok: boolean;
  result?: any;
  errors?: string[];
  error?: string;
}


