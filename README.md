# Simplified Twitter Client-Server System in Go

This project implements a simplified Twitter client-server system to practice low-level parallel programming and data structure design. It includes features such as a thread-safe singly linked list, a client-server communication model, and benchmarking for scalability and performance.

## Features
- Core Twitter feed operations implemented using a singly linked list.
- Thread-safe feed with a custom read-write lock supporting up to 32 concurrent readers.
- Client-server interaction using JSON-encoded requests and responses.
- Parallel task processing with a lock-free queue using the producer-consumer model.

## Project Structure
```
simplified-twitter/
├── feed.go               # Twitter feed implementation with thread safety
├── server.go             # Client-server interaction and parallel task processing
├── twitter.go            # Main entry point and benchmarking logic
├── benchmark/            # Benchmarking scripts and output analysis
```

### Key Files
- **`feed.go`**: Manages feed operations and thread safety using a custom RW lock.
- **`server.go`**: Implements a client-server model with sequential and parallel modes.
- **`twitter.go`**: Entry point for configuring and running the program.

## Installation

### Prerequisites
- **Golang 1.20** or later.

### Steps
1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/simplified-twitter.git
   cd simplified-twitter/simplified-twitter
   ```

## Usage
Run benchmarks using the provided script:
```bash
cd benchmark
sbatch run.sh
```
The script generates:
- **Benchmark directories** (`xsmall`, `small`, `medium`, `large`, `xlarge`) containing runtime data.
- **`speedup.png`**: Speedup graph comparing parallel and sequential implementations.

## Performance Analysis

For detailed performance analysis and results, refer to the [Performance Analysis](benchmark/performance_analysis.md) document.