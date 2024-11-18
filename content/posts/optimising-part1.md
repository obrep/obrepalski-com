+++
title = 'Optimising Cloud Services At Scale'
date = 2024-11-17
tags = ["series", "docs"]
series =  ["Optimising Cloud Services At Scale"]
draft = true
series_order =  2
+++
# Intro


This article series is a guide to optimising your Cloud services at scale. While parts of it will be specific to technologies I'm using daily (Go, GCP), most of the practices here will be universal for most of the services running on Kubernetes, which has become the de facto standard for large Deployments over the past decade.

We will start by looking into tools and techniques for optimising the service and in later parts expand the scope to include infrastructure, networking etc.

Topics covered:
- Service-level optimisations: profiling and tuning
- Continuous profiling
- Measuring impact
- Choosing right infrastructure
- Right-sizing
- Autoscaling
- What else am I paying for?
- Cheaper compute: commited use and spot
- Networking optimisations
- Closing thought and my experiences

When working with microservices at scale, performance optimization becomes crucial for maintaining both system reliability and cost efficiency. While this post covers concepts applicable to any language, the tools and infrastructure discussed focus on Go and Google Cloud Platform, as that's where my experience lies. Most services my team maintains have multiple downstream dependencies, making traditional end-to-end performance testing challenging.

## Why Not Just Load Test?

TODO: rethink this section

Load testing is perfect for certain scenarios - like testing a stateless API service, a data processing pipeline, or a standalone web application. However, when dealing with services that have multiple downstream dependencies running in production, load testing becomes significantly more complex. Setting up a realistic load test would require either a sophisticated staging environment mirroring production or careful coordination to avoid impacting production dependencies. This is where profiling comes in as a powerful alternative.

# Optimising Go Services

Go has a few parameters that we can tune to  

## Local Profiling with pprof / Understanding what our service is doing

### Setup
If you want to follow along, you will need to have Go installed 

Go's built-in `pprof` tool provides the foundation for most profiling work:
### 
```go
package main 

import (
    "net/http"
    _ "net/http/pprof"
    "log"
)

func main() {
    // Add debug endpoint for pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // Your service code here
}
```

This will expose several profiling endpoints:
- `/debug/pprof/profile` - CPU profile
- `/debug/pprof/heap` - Memory profile
- `/debug/pprof/goroutine` - Goroutine profile
- `/debug/pprof/block` - Block profile
- `/debug/pprof/mutex` - Mutex profile

The profile can then be read using
```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

```
This will produce a `profile.pb`


### Understanding Flame Graphs

Flame graphs visualize profiling data in a hierarchical format, making it easier to spot performance bottlenecks. The width of each function block represents the time spent in that function, while the height shows the call stack depth. Functions consuming more CPU time appear wider, making them easy to identify.

You will need to install graphviz to open http, depending on your system it may be done by running the following commands or [downloading](https://graphviz.org/download/):
```bash
brew install graphviz # MacOs
apt install graphviz # Ubuntu/Debian
```

```bash
# Interactive web UI - most useful for me 
go tool pprof -http=:8080 [profile.pb.gz]
```


[INSERT FLAME GRAPH SCREENSHOT HERE - Example of a CPU profile from one of our services showing a typical pattern with GC overhead]

## Profiling Fundamentals

Profiling offers insights into runtime behavior, helping identify bottlenecks and performance issues that might not be apparent through code review alone. The two main types of profiling most useful in production environments are:

**CPU Profiling** captures where programs spend their processing time, identifying hot paths and expensive operations in the code.

**Memory Profiling** reveals memory usage patterns, helping understand heap allocations and garbage collection behavior. This becomes particularly important in Go, where understanding allocation patterns can significantly impact performance.

Go also provides specialized profiling for goroutines, blocking operations, and mutex contention, though these are typically less frequently used in day-to-day optimization work.

### Tuning Go's Garbage Collector

Go provides two main knobs for GC tuning:

**GOGC**: Controls the garbage collector's target heap size relative to live heap. The default value of 100 means the heap can grow to 100% larger than the live heap before triggering a collection.

**GOMEMLIMIT** (Go 1.19+): Sets a hard memory limit for the heap. Based on personal experience, disabling GOGC (`GOGC=off`) and using GOMEMLIMIT alone often provides good results. However, when the live heap approaches the set limit, the application can experience thrashing as the garbage collector works constantly to keep memory usage below the limit. This can lead to significant CPU overhead and reduced performance. Setting GOMEMLIMIT about 20-30% higher than expected peak live heap usage usually provides a good balance.

### CPU Management with GOMAXPROCS

GOMAXPROCS controls the number of operating system threads that can execute Go code simultaneously. While Go usually manages this well automatically, there are cases where adjustment helps:
- When running in containerized environments with CPU limits (note that Go will detect the number of cores on the underlying node instead of the container's CPU limit)
- When dealing with heavy I/O operations
- When CPU quotas don't align with available cores

### Additional Metrics with Prometheus

Beyond basic profiling, exposing metrics to Prometheus provides valuable insights into application behavior. The Go Prometheus client includes collectors for runtime metrics:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/collectors"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    reg := prometheus.NewRegistry()
    
    // Add Go collectors
    reg.MustRegister(
        collectors.NewGoCollector(),
        collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
    )

    // Expose metrics
    http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
}
```

For a complete example of runtime metrics collection, see the [prometheus/client_golang example](https://github.com/prometheus/client_golang/blob/main/examples/gocollector/main.go).

## Continuous Profiling

While point-in-time profiling is useful during development, continuous profiling in production provides insights into real-world performance patterns.

### Google Cloud Profiler

Google Cloud Profiler can be easily integrated into Go services:

```go
import (
    "cloud.google.com/go/profiler"
    "os"
)

func main() {
    if err := profiler.Start(profiler.Config{
        Service:        "myservice",
        ServiceVersion: os.Getenv("BUILD_ID"),  // [1]
    }); err != nil {
        log.Fatal(err)
    }
}
```
[1] Using BUILD_ID from the CI/CD system allows tracking performance changes across deployments and observing the history of specific function calls by version

Note: The service account used for deployment needs the `roles/cloudprofiler.agent` IAM role to submit profiles.

For detailed setup and usage, refer to the [Google Cloud Profiler documentation](https://cloud.google.com/profiler/docs/).

### Other Profiling Solutions

Major cloud providers offer their own profiling solutions, though they typically focus on specific languages - AWS CodeGuru targets Java and Python, while Azure Monitor primarily supports .NET applications. Third-party solutions like [Datadog](https://docs.datadoghq.com/profiler/), [New Relic](https://newrelic.com/platform/application-monitoring), [Honeycomb](https://www.honeycomb.io/), and [Grafana Pyroscope](https://grafana.com/oss/pyroscope/) provide language-agnostic profiling capabilities integrated with their observability platforms.

## Measuring Impact

Establishing a baseline before optimization is crucial. Cloud Profiler makes this particularly easy by enabling version-to-version comparisons. A practical approach is comparing performance between canary deployments running new code and current production instances.

When measuring the impact of optimizations, several key metrics should be monitored as guardrails:
- Request latency (p50 and p99)
- Error rates, particularly 5xx responses
- Pod restarts and OOM events
- Resource utilization (CPU, memory)
- Garbage collection metrics

These metrics help ensure that optimizations don't inadvertently cause regressions in other areas of the service.

## Further Reading

For deeper understanding of profiling and garbage collection in Go:
- [Profiling Go Programs](https://go.dev/blog/pprof) - Official Go blog post on pprof
- [Go Garbage Collector Guide](https://tip.golang.org/doc/gc-guide) - Comprehensive guide to Go's garbage collector

[This article will continue with specific optimization techniques and case studies in part 2]
