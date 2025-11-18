/**
 * Example usage of structured-parse TypeScript library
 */

import { createParser, Label } from './src/index.js';

async function main() {
  console.log('🚀 structured-parse TypeScript Example\n');

  // Create and initialize the parser
  console.log('Initializing parser (loading WASM)...');
  const parser = await createParser();
  console.log('✅ Parser ready!\n');

  // Example 1: Basic parsing with JSON field
  console.log('Example 1: Basic Parse');
  console.log('─'.repeat(50));

  const labels1: Label[] = [
    { name: 'Action', required: true },
    { name: 'Parameters', isJson: true },
    { name: 'Thought' }
  ];

  const text1 = `
Action: process_data
Parameters: {"input": ["file1.txt", "file2.txt"], "mode": "batch"}
Thought: Processing multiple files in batch mode
  `.trim();

  const result1 = parser.parse(labels1, text1);

  console.log('Input:', text1);
  console.log('\nParsed Result:');
  console.log(JSON.stringify(result1.result, null, 2));
  console.log('Errors:', result1.errors.length === 0 ? 'None' : result1.errors);

  // Example 2: Block parsing
  console.log('\n\nExample 2: Block Parsing');
  console.log('─'.repeat(50));

  const labels2: Label[] = [
    { name: 'Task', isBlockStart: true },
    { name: 'Input', isJson: true },
    { name: 'Result' }
  ];

  const text2 = `
Task: Summarize
Input: {"text": "Long article about AI..."}
Result: AI article discusses recent developments

Task: Classify
Input: {"text": "Product review..."}
Result: Positive sentiment
  `.trim();

  const result2 = parser.parseBlocks(labels2, text2);

  console.log('Input:', text2);
  console.log('\nParsed Blocks:');
  result2.blocks.forEach((block, i) => {
    console.log(`\nBlock ${i + 1}:`, JSON.stringify(block, null, 2));
  });
  console.log('Errors:', result2.errors.length === 0 ? 'None' : result2.errors);

  // Example 3: Validation (required field missing)
  console.log('\n\nExample 3: Validation Errors');
  console.log('─'.repeat(50));

  const labels3: Label[] = [
    { name: 'Name', required: true },
    { name: 'Email', requiredWith: ['Name'] }
  ];

  const text3 = 'Email: test@example.com';  // Missing required 'Name'

  const result3 = parser.parse(labels3, text3);

  console.log('Input:', text3);
  console.log('\nParsed Result:', JSON.stringify(result3.result, null, 2));
  console.log('Errors:', result3.errors);

  console.log('\n✅ Examples complete!');
}

main().catch(error => {
  console.error('❌ Error:', error);
  process.exit(1);
});


