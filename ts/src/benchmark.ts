import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';
import { createParser } from './index.js';
import type { Label } from './types.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const assetsDir = path.join(__dirname, '../../test-assets');

interface BenchmarkResult {
  name: string;
  iterations: number;
  totalTimeMs: number;
  avgTimeMs: number;
  opsPerSec: number;
}

async function benchmark(name: string, fn: () => void, iterations: number = 1000): Promise<BenchmarkResult> {
  // Warm up
  for (let i = 0; i < 10; i++) {
    fn();
  }
  
  // Measure
  const start = performance.now();
  for (let i = 0; i < iterations; i++) {
    fn();
  }
  const end = performance.now();
  
  const totalTimeMs = end - start;
  const avgTimeMs = totalTimeMs / iterations;
  const opsPerSec = 1000 / avgTimeMs;
  
  return {
    name,
    iterations,
    totalTimeMs,
    avgTimeMs,
    opsPerSec,
  };
}

async function runBenchmarks() {
  console.log('🚀 Starting TypeScript/WASM Benchmarks\n');
  console.log('═'.repeat(80));
  
  const parser = await createParser();
  const results: BenchmarkResult[] = [];
  
  // Benchmark 1: Basic functionality
  {
    const input = fs.readFileSync(path.join(assetsDir, 'basic_functionality_input.txt'), 'utf-8');
    const labels: Label[] = [
      { name: 'Action Input', requiredWith: ['Action'], isJson: true },
      { name: 'Action', requiredWith: ['Action Input'] },
      { name: 'Thought' },
      { name: 'Result', required: true },
    ];
    
    const result = await benchmark('Basic Parse', () => {
      parser.parse(labels, input);
    }, 1000);
    results.push(result);
  }
  
  // Benchmark 2: Mixed case multiline
  {
    const input = fs.readFileSync(path.join(assetsDir, 'mixed_case_multiline_input.txt'), 'utf-8');
    const labels: Label[] = [
      { name: 'Context' },
      { name: 'Intention' },
      { name: 'Role' },
      { name: 'Action' },
      { name: 'Outcome' },
      { name: 'Notes' },
    ];
    
    const result = await benchmark('Multiline Parse', () => {
      parser.parse(labels, input);
    }, 1000);
    results.push(result);
  }
  
  // Benchmark 3: JSON parsing
  {
    const input = fs.readFileSync(path.join(assetsDir, 'json_and_malformed_input.txt'), 'utf-8');
    const labels: Label[] = [
      { name: 'Config', isJson: true },
      { name: 'Data', isJson: true },
      { name: 'Description' },
    ];
    
    const result = await benchmark('JSON Parse', () => {
      parser.parse(labels, input);
    }, 1000);
    results.push(result);
  }
  
  // Benchmark 4: Block parsing
  {
    const input = fs.readFileSync(path.join(assetsDir, 'block_parsing_input.txt'), 'utf-8');
    const labels: Label[] = [
      { name: 'Task', isBlockStart: true },
      { name: 'Input', isJson: true },
      { name: 'Result' },
    ];
    
    const result = await benchmark('Block Parse', () => {
      parser.parseBlocks(labels, input);
    }, 1000);
    results.push(result);
  }
  
  // Print results
  console.log('\nBenchmark Results:');
  console.log('─'.repeat(80));
  console.log('Test Name'.padEnd(25), 'Iterations'.padEnd(12), 'Avg Time'.padEnd(15), 'Ops/sec');
  console.log('─'.repeat(80));
  
  for (const result of results) {
    console.log(
      result.name.padEnd(25),
      result.iterations.toString().padEnd(12),
      `${result.avgTimeMs.toFixed(3)} ms`.padEnd(15),
      result.opsPerSec.toFixed(2)
    );
  }
  
  console.log('═'.repeat(80));
  
  // Calculate overall stats
  const totalOps = results.reduce((sum, r) => sum + r.iterations, 0);
  const totalTime = results.reduce((sum, r) => sum + r.totalTimeMs, 0);
  const avgOpsPerSec = results.reduce((sum, r) => sum + r.opsPerSec, 0) / results.length;
  
  console.log('\nOverall Statistics:');
  console.log(`  Total operations: ${totalOps}`);
  console.log(`  Total time: ${totalTime.toFixed(2)} ms`);
  console.log(`  Average ops/sec: ${avgOpsPerSec.toFixed(2)}`);
  console.log('\n✅ Benchmarks complete!\n');
}

runBenchmarks().catch(error => {
  console.error('❌ Benchmark failed:', error);
  process.exit(1);
});


