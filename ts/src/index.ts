import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';
import type { Label, ParserOptions, ParseResult, ParseBlocksResult, WasmResponse } from './types.js';

export type { Label, ParserOptions, ParseResult, ParseBlocksResult };

/**
 * StructuredParser provides WebAssembly-powered parsing of labeled/structured LLM output.
 */
export class StructuredParser {
  private wasmReady: Promise<void>;
  private wasmParse: ((request: string) => string) | null = null;
  private wasmParseBlocks: ((request: string) => string) | null = null;
  private wasmVersion: (() => string) | null = null;

  /**
   * Creates a new StructuredParser instance.
   * The WASM module is loaded asynchronously; use await parser.ready() before parsing.
   * 
   * @param wasmPath - Optional path to the WASM file. If not provided, uses the bundled WASM module.
   */
  constructor(wasmPath?: string) {
    this.wasmReady = this.init(wasmPath);
  }

  /**
   * Waits for the WASM module to be fully loaded and ready.
   */
  async ready(): Promise<void> {
    await this.wasmReady;
  }

  /**
   * Returns the version of the WASM module.
   */
  version(): string {
    if (!this.wasmVersion) {
      throw new Error('WASM module not initialized. Call await parser.ready() first.');
    }
    return this.wasmVersion();
  }

  /**
   * Parses text into a map of label names to their values.
   * 
   * @param labels - Array of label definitions
   * @param text - The text to parse
   * @param options - Optional parser options
   * @returns ParseResult containing the parsed data and any errors
   */
  parse(labels: Label[], text: string, options?: ParserOptions): ParseResult {
    if (!this.wasmParse) {
      throw new Error('WASM module not initialized. Call await parser.ready() first.');
    }

    const request = {
      labels,
      text,
      options: options || undefined,
    };

    const responseJSON = this.wasmParse(JSON.stringify(request));
    const response: WasmResponse = JSON.parse(responseJSON);

    if (!response.ok && response.error) {
      throw new Error(`Parse failed: ${response.error}`);
    }

    return {
      result: response.result || {},
      errors: response.errors || [],
    };
  }

  /**
   * Parses text into blocks, splitting at the block start label.
   * 
   * @param labels - Array of label definitions (one must have isBlockStart: true)
   * @param text - The text to parse
   * @param options - Optional parser options
   * @returns ParseBlocksResult containing the parsed blocks and any errors
   */
  parseBlocks(labels: Label[], text: string, options?: ParserOptions): ParseBlocksResult {
    if (!this.wasmParseBlocks) {
      throw new Error('WASM module not initialized. Call await parser.ready() first.');
    }

    const request = {
      labels,
      text,
      options: options || undefined,
    };

    const responseJSON = this.wasmParseBlocks(JSON.stringify(request));
    const response: WasmResponse = JSON.parse(responseJSON);

    if (!response.ok && response.error) {
      throw new Error(`ParseBlocks failed: ${response.error}`);
    }

    return {
      blocks: response.result || [],
      errors: response.errors || [],
    };
  }

  /**
   * Initializes the WASM module.
   * @private
   */
  private async init(wasmPath?: string): Promise<void> {
    // Determine the path to the WASM file
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = path.dirname(__filename);
    const defaultWasmPath = path.join(__dirname, 'wasm', 'structured-parse.wasm');
    const wasmFilePath = wasmPath || defaultWasmPath;

    // Load the Go WASM exec helper
    const wasmExecPath = path.join(__dirname, 'wasm', 'wasm_exec.js');
    
    // Dynamically import the wasm_exec.js which sets up the global Go object
    await import(wasmExecPath);

    // Load the WASM module
    const wasmBuffer = fs.readFileSync(wasmFilePath);
    
    // @ts-ignore - Go is added by wasm_exec.js
    const go = new Go();
    const result = await WebAssembly.instantiate(wasmBuffer, go.importObject);
    
    // Run the Go WASM module (this sets up the exported functions)
    go.run(result.instance);

    // Wait a bit for the Go runtime to initialize
    await new Promise(resolve => setTimeout(resolve, 100));

    // Get references to the exported functions
    // @ts-ignore - these are set by the Go WASM module
    this.wasmParse = globalThis.wasmParse;
    // @ts-ignore
    this.wasmParseBlocks = globalThis.wasmParseBlocks;
    // @ts-ignore
    this.wasmVersion = globalThis.wasmVersion;

    if (!this.wasmParse || !this.wasmParseBlocks || !this.wasmVersion) {
      throw new Error('Failed to load WASM exports. Make sure the WASM module is built correctly.');
    }
  }
}

/**
 * Creates a new parser instance and waits for it to be ready.
 * This is a convenience function that combines construction and initialization.
 * 
 * @param wasmPath - Optional path to the WASM file
 * @returns A ready-to-use StructuredParser instance
 */
export async function createParser(wasmPath?: string): Promise<StructuredParser> {
  const parser = new StructuredParser(wasmPath);
  await parser.ready();
  return parser;
}


