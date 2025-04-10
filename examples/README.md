# Examples from obrepalski.com

This repository contains examples from [obrepalski.com](https://obrepalski.com) blog posts, designed to be accessible for people who are new to Go programming.

## Getting Started with Go

If you're new to Go, here's how to get started:

1. **Install Go**: Download and install Go from [golang.org](https://golang.org/dl/)
2. **Verify installation**: Run `go version` in your terminal to ensure Go is properly installed
3. **Download dependencies**:  Run `go mod download` in `examples/` folder. The `go.mod` and `go.sum` files specify the dependencies (and will be used when running the command)

## How to Run Examples

1. Ensure Go is installed on your system
2. Clone this repository
    ```
    go mod download
    ```
3. Navigate to an example directory
4. Run the example with:
    ```
    go run .
    ```
5. For benchmarks, use:
    ```
    go test -bench=. .
    ```

### Useful Commands for Performance Tuning

1. **Obtain CPU profile:** (assumes your application is already running and you've exposed pprof on port 6060)
    ```bash
    go tool pprof -http=:9090 localhost:6060/debug/pprof/profile
    ```

2. **Obtain memory profile:**
    ```bash
    go tool pprof -http=:9090 localhost:6060/debug/pprof/heap
    ```

3. **Run benchmarks with memory allocation stats:**
    ```bash
    go test -bench=. -benchmem
    ```

4. **Generate a CPU profile during benchmarks:**
    ```bash
    go test -bench=. -cpuprofile=cpu.prof
    ```

5. **Generate a memory profile during benchmarks:**
    ```bash
    go test -bench=. -memprofile=mem.prof
    ```

6. **Visualize execution traces:**
    ```bash
    go test -trace=trace.out
    go tool trace trace.out
    ```

These commands are useful for profiling and benchmarking your Go applications, as well as for debugging performance issues.

## Additional Learning Resources

- [obrepalski.com](https://obrepalski.com) - Posts with context for the examples
- [A Tour of Go](https://tour.golang.org/) - Interactive tutorials for getting started with Go. My go to recommendation for onboarding, especially for people with experience in other languages.
- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)