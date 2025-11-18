"""Benchmark script for structured-parse Python bindings."""

import json
import time
from pathlib import Path
from dataclasses import dataclass
from typing import Callable

from structured_parse import StructuredParser, Label

# Get the path to test assets
ASSETS_DIR = Path(__file__).parent.parent / "test-assets"


@dataclass
class BenchmarkResult:
    """Result of a benchmark run."""
    name: str
    iterations: int
    total_time_ms: float
    avg_time_ms: float
    ops_per_sec: float


def benchmark(name: str, fn: Callable, iterations: int = 1000) -> BenchmarkResult:
    """
    Run a benchmark for the given function.
    
    Args:
        name: Name of the benchmark
        fn: Function to benchmark
        iterations: Number of iterations to run
        
    Returns:
        BenchmarkResult with timing statistics
    """
    # Warm up
    for _ in range(10):
        fn()
    
    # Measure
    start = time.perf_counter()
    for _ in range(iterations):
        fn()
    end = time.perf_counter()
    
    total_time_ms = (end - start) * 1000
    avg_time_ms = total_time_ms / iterations
    ops_per_sec = 1000 / avg_time_ms
    
    return BenchmarkResult(
        name=name,
        iterations=iterations,
        total_time_ms=total_time_ms,
        avg_time_ms=avg_time_ms,
        ops_per_sec=ops_per_sec,
    )


def run_benchmarks():
    """Run all benchmarks."""
    print("🚀 Starting Python/WASM Benchmarks\n")
    print("=" * 80)
    
    parser = StructuredParser()
    results = []
    
    # Benchmark 1: Basic functionality
    with open(ASSETS_DIR / "basic_functionality_input.txt") as f:
        input_text = f.read()
    labels1 = [
        Label(name="Action Input", required_with=["Action"], is_json=True),
        Label(name="Action", required_with=["Action Input"]),
        Label(name="Thought"),
        Label(name="Result", required=True),
    ]
    
    result = benchmark("Basic Parse", lambda: parser.parse(labels1, input_text), 1000)
    results.append(result)
    
    # Benchmark 2: Mixed case multiline
    with open(ASSETS_DIR / "mixed_case_multiline_input.txt") as f:
        input_text = f.read()
    labels2 = [
        Label(name="Context"),
        Label(name="Intention"),
        Label(name="Role"),
        Label(name="Action"),
        Label(name="Outcome"),
        Label(name="Notes"),
    ]
    
    result = benchmark("Multiline Parse", lambda: parser.parse(labels2, input_text), 1000)
    results.append(result)
    
    # Benchmark 3: JSON parsing
    with open(ASSETS_DIR / "json_and_malformed_input.txt") as f:
        input_text = f.read()
    labels3 = [
        Label(name="Config", is_json=True),
        Label(name="Data", is_json=True),
        Label(name="Description"),
    ]
    
    result = benchmark("JSON Parse", lambda: parser.parse(labels3, input_text), 1000)
    results.append(result)
    
    # Benchmark 4: Block parsing
    with open(ASSETS_DIR / "block_parsing_input.txt") as f:
        input_text = f.read()
    labels4 = [
        Label(name="Task", is_block_start=True),
        Label(name="Input", is_json=True),
        Label(name="Result"),
    ]
    
    result = benchmark("Block Parse", lambda: parser.parse_blocks(labels4, input_text), 1000)
    results.append(result)
    
    # Print results
    print("\nBenchmark Results:")
    print("-" * 80)
    print(f"{'Test Name':<25} {'Iterations':<12} {'Avg Time':<15} {'Ops/sec'}")
    print("-" * 80)
    
    for result in results:
        print(
            f"{result.name:<25} "
            f"{result.iterations:<12} "
            f"{result.avg_time_ms:>8.3f} ms     "
            f"{result.ops_per_sec:>8.2f}"
        )
    
    print("=" * 80)
    
    # Calculate overall stats
    total_ops = sum(r.iterations for r in results)
    total_time = sum(r.total_time_ms for r in results)
    avg_ops_per_sec = sum(r.ops_per_sec for r in results) / len(results)
    
    print("\nOverall Statistics:")
    print(f"  Total operations: {total_ops}")
    print(f"  Total time: {total_time:.2f} ms")
    print(f"  Average ops/sec: {avg_ops_per_sec:.2f}")
    print("\n✅ Benchmarks complete!\n")


if __name__ == "__main__":
    run_benchmarks()


