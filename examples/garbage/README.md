# Garbage Collection Demo Application

This is an intentionally unoptimized Go application designed to demonstrate the impact of garbage collection on application performance through flamegraphs and profiling.

## Purpose

This application showcases various anti-patterns and inefficiencies to help visualize:
- Impact of excessive heap allocations on GC
- Effects of GOGC and GOMEMLIMIT tuning
- Memory pressure from poor coding practices
- CPU overhead from frequent garbage collection

## Anti-Patterns Implemented

### 1. Poorly Aligned Structs
- `User` and `Transaction` structs with inefficient field ordering
- Causes memory padding and wastes ~40% of struct space

### 2. Inefficient String Operations
- String concatenation using `+` operator in loops
- Creates many temporary string allocations
- Should use `strings.Builder` instead

### 3. Excessive Heap Allocations
- Interface{} boxing causing heap escapes
- Temporary slice allocations in loops
- Inefficient slice trimming (creating new backing arrays)
- Map recreation instead of deletion

### 4. Memory Leaks
- `leakyGoroutine()` accumulates data indefinitely
- Session data never cleaned up
- Simulates common memory leak patterns

### 5. Inefficient Data Structures
- Nested maps (`map[int]map[string]string`)
- Using `interface{}` for caching (causes boxing)
- Copying entire slices unnecessarily

## Running the Application

### Basic Run
```bash
go run main.go
```

### With GC Tuning
```bash
# Disable GC (not recommended for production!)
GOGC=off go run main.go

# Increase GC threshold (less frequent GC)
GOGC=200 go run main.go

# Set memory limit (Go 1.19+)
GOMEMLIMIT=500MiB go run main.go

# Combine settings
GOGC=off GOMEMLIMIT=1GiB go run main.go
```

### With Profiling Tools
```bash
# Enable GC trace
GODEBUG=gctrace=1 go run main.go

# Enable memory profiler
GODEBUG=gctrace=1,gcpacertrace=1 go run main.go
```

## Generating Flamegraphs

### 1. Start the application
```bash
go run main.go
```

### 2. Capture CPU Profile
```bash
# 30-second CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Or save directly
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof
```

### 3. Capture Heap Profile
```bash
# Current heap usage
go tool pprof http://localhost:6060/debug/pprof/heap

# Allocation profile
go tool pprof http://localhost:6060/debug/pprof/allocs
```

### 4. Generate Flamegraph
```bash
# Using pprof (interactive)
go tool pprof -http=:8080 cpu.prof

# Export to SVG
go tool pprof -svg cpu.prof > cpu_flamegraph.svg

# Using FlameGraph tools
go tool pprof -raw cpu.prof | FlameGraph/stackcollapse-go.pl | FlameGraph/flamegraph.pl > flamegraph.svg
```

## Metrics to Observe

The application prints memory statistics every 10 seconds:
- **Alloc**: Current heap allocation
- **TotalAlloc**: Cumulative allocation
- **NumGC**: Number of GC cycles
- **HeapObjects**: Number of allocated objects

## Optimization Opportunities

Future improvements to demonstrate:

1. **Struct Alignment**: Reorder fields by size
2. **String Building**: Use `strings.Builder`
3. **Object Pooling**: Implement `sync.Pool`
4. **Buffer Reuse**: Reuse slices and buffers
5. **Reduce Boxing**: Use concrete types instead of `interface{}`
6. **Fix Memory Leaks**: Implement proper cleanup
7. **Efficient Data Structures**: Use appropriate containers

## Expected Results

### Unoptimized (Default)
- High GC frequency (700+ GCs in 10 seconds)
- Significant CPU time in `runtime.gcBgMarkWorker`
- Memory constantly growing and shrinking
- High allocation rate (~90 MB/s)

### With GOGC=200
- Reduced GC frequency (~50% less)
- Higher memory usage
- Better throughput

### With GOMEMLIMIT=500MiB
- Predictable memory usage
- Potential GC thrashing near limit
- More consistent performance

### With Optimizations
- 80-90% reduction in allocations
- 10x reduction in GC frequency
- Significantly improved CPU efficiency
- Stable memory usage

## Monitoring URLs

When running, access these endpoints:
- http://localhost:6060/debug/pprof/ - Profile index
- http://localhost:6060/debug/pprof/heap - Heap profile
- http://localhost:6060/debug/pprof/profile?seconds=30 - CPU profile
- http://localhost:6060/debug/pprof/allocs - Allocation profile
- http://localhost:6060/debug/pprof/goroutine - Goroutine profile