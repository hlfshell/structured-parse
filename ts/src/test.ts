import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';
import { createParser } from './index.js';
import type { Label } from './types.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const assetsDir = path.join(__dirname, '../../test-assets');

// Test helper
function assert(condition: boolean, message: string): void {
  if (!condition) {
    throw new Error(`Assertion failed: ${message}`);
  }
}

function deepEqual(a: any, b: any): boolean {
  if (a === b) return true;
  if (typeof a !== typeof b) return false;
  if (typeof a !== 'object' || a === null || b === null) return false;
  
  const keysA = Object.keys(a);
  const keysB = Object.keys(b);
  if (keysA.length !== keysB.length) return false;
  
  for (const key of keysA) {
    if (!keysB.includes(key)) return false;
    if (!deepEqual(a[key], b[key])) return false;
  }
  
  return true;
}

async function testBasicFunctionality() {
  console.log('Testing basic functionality...');
  
  const input = fs.readFileSync(path.join(assetsDir, 'basic_functionality_input.txt'), 'utf-8');
  const expected = JSON.parse(fs.readFileSync(path.join(assetsDir, 'basic_functionality_output.json'), 'utf-8'));
  
  const labels: Label[] = [
    { name: 'Action Input', requiredWith: ['Action'], isJson: true },
    { name: 'Action', requiredWith: ['Action Input'] },
    { name: 'Thought' },
    { name: 'Result', required: true },
  ];
  
  const parser = await createParser();
  const result = parser.parse(labels, input);
  
  assert(result.errors.length === 0, `Expected no errors, got: ${result.errors.join(', ')}`);
  assert(deepEqual(result.result, expected), `Result mismatch:\nGot: ${JSON.stringify(result.result, null, 2)}\nExpected: ${JSON.stringify(expected, null, 2)}`);
  
  console.log('✅ Basic functionality test passed');
}

async function testMixedCaseMultiline() {
  console.log('Testing mixed case and multiline...');
  
  const input = fs.readFileSync(path.join(assetsDir, 'mixed_case_multiline_input.txt'), 'utf-8');
  const expected = JSON.parse(fs.readFileSync(path.join(assetsDir, 'mixed_case_multiline_output.json'), 'utf-8'));
  
  const labels: Label[] = [
    { name: 'Context' },
    { name: 'Intention' },
    { name: 'Role' },
    { name: 'Action' },
    { name: 'Outcome' },
    { name: 'Notes' },
  ];
  
  const parser = await createParser();
  const result = parser.parse(labels, input);
  
  assert(result.errors.length === 0, `Expected no errors, got: ${result.errors.join(', ')}`);
  assert(deepEqual(result.result, expected), `Result mismatch:\nGot: ${JSON.stringify(result.result, null, 2)}\nExpected: ${JSON.stringify(expected, null, 2)}`);
  
  console.log('✅ Mixed case multiline test passed');
}

async function testJSONAndMalformed() {
  console.log('Testing JSON and malformed JSON...');
  
  const input = fs.readFileSync(path.join(assetsDir, 'json_and_malformed_input.txt'), 'utf-8');
  const expected = JSON.parse(fs.readFileSync(path.join(assetsDir, 'json_and_malformed_output.json'), 'utf-8'));
  const expectedErrors = JSON.parse(fs.readFileSync(path.join(assetsDir, 'json_and_malformed_errors.json'), 'utf-8'));
  
  const labels: Label[] = [
    { name: 'Config', isJson: true },
    { name: 'Data', isJson: true },
    { name: 'Description' },
  ];
  
  const parser = await createParser();
  const result = parser.parse(labels, input);
  
  assert(deepEqual(result.result, expected), `Result mismatch:\nGot: ${JSON.stringify(result.result, null, 2)}\nExpected: ${JSON.stringify(expected, null, 2)}`);
  assert(result.errors.length === expectedErrors.length, `Expected ${expectedErrors.length} errors, got ${result.errors.length}`);
  
  console.log('✅ JSON and malformed test passed');
}

async function testBlockParsing() {
  console.log('Testing block parsing...');
  
  const input = fs.readFileSync(path.join(assetsDir, 'block_parsing_input.txt'), 'utf-8');
  const expected = JSON.parse(fs.readFileSync(path.join(assetsDir, 'block_parsing_output.json'), 'utf-8'));
  
  const labels: Label[] = [
    { name: 'Task', isBlockStart: true },
    { name: 'Input', isJson: true },
    { name: 'Result' },
  ];
  
  const parser = await createParser();
  const result = parser.parseBlocks(labels, input);
  
  assert(result.errors.length === 0, `Expected no errors, got: ${result.errors.join(', ')}`);
  assert(deepEqual(result.blocks, expected), `Blocks mismatch:\nGot: ${JSON.stringify(result.blocks, null, 2)}\nExpected: ${JSON.stringify(expected, null, 2)}`);
  
  console.log('✅ Block parsing test passed');
}

async function runAllTests() {
  try {
    console.log('🧪 Running TypeScript tests...\n');
    
    await testBasicFunctionality();
    await testMixedCaseMultiline();
    await testJSONAndMalformed();
    await testBlockParsing();
    
    console.log('\n✅ All TypeScript tests passed!');
  } catch (error) {
    console.error('\n❌ Test failed:', error);
    process.exit(1);
  }
}

runAllTests();


