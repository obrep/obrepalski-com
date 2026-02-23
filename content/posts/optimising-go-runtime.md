+++
title = 'Go Runtime'
date = 2024-01-30
tags = ["series", "docs"]
series =  ["Optimising Cloud Services At Scale"]
draft = true
series_order =  4
+++
# Go Runtime: Brief Overview

Before tuning the runtime we will start with a high level overview of two key Go components: scheduler and garbage collector

## Garbage collection

Go is garbage collected language.  -> objects allocated on heap need to be periodically cleaned up
Memory allocated on stack will automatically be freed up when the function returns
Can decrease how often GC needs to run by decreasing allocations
Not always clear if a variable will be allocated on heap
Check where we allocate most using the heap prole, similarly to cpu prole (link to previous article)

Nice visualization: https://pusher.github.io/tricolor-gc-visualization/

## Scheduler 



# Tuning

Go has a few parameters that we can tune to achieve better performance.


## Tuning Go's Garbage Collector - `GOGC` and `GOMEMLIMIT`

Go provides two main knobs for GC tuning:

**GOGC**: Controls the garbage collector's target heap size relative to live heap. The default value of 100 means the heap can grow to 100% larger than the live heap before triggering a collection.

**GOMEMLIMIT** (Go 1.19+): Sets a hard memory limit for the heap. Based on personal experience, disabling GOGC (`GOGC=off`) and using GOMEMLIMIT alone often provides good results. However, when the live heap approaches the set limit, the application can experience thrashing as the garbage collector works constantly to keep memory usage below the limit. This can lead to significant CPU overhead and reduced performance. Setting GOMEMLIMIT about 20-30% higher than expected peak live heap usage usually provides a good balance.

Go provides a way to specify the cpu/memory tradeoff that is best for our scenario
Default behaviour is to wait until heaps grows by 100%
GOGC (default=100 ) - how much live heap can grow before triggering GC
GOMEMLIMIT (default=off ) - soft memory limit (from Go 1.19)
Possible to still OOM with GOMEMLIMIT set (trashing)
Default settings are quite sensible but adjusting them can give us some extra performance


Tips:
Look for runtime.gc* in CPU prole
Set GODEBUG=gctrace=1 for detailed trace whenever collection happens
Start with GOMEMLIMIT at ~90% of available memory
Take into account different price of resources (each core costs close to 7GB memory)

For a more in-depth guide on Go's GC, refer to the official documentation.

Adjusting those values is as simple as setting environment variables

``shell
GOGC="off" go run .
``

# Disable garbage collector
GOGC="off" go run .

# Kubernetes configuration example
```yaml
spec:
  containers:
  - name: your-app
    image: your-image
    env:
    - name: GOMAXPROCS
      value: "10"
    - name: GOGC
      value: "off"
    - name: GOMEMLIMIT
      value: 2500
```


## CPU Management with GOMAXPROCS

Denes how many threads will be used for simultaneously running service code
Value too low -> you won’t utilise all resources
Value too high -> may lead to longer scheduling latency and high tail latencies
Best place to start is to have GOMAXPROCS equal to number of available cores (which is also the default
behaviour)
High values when running on kubernetes as it will default to the # of cores of the underlying instance
Two possible workarounds:
Use automaxprocs
Set the env values based on requests/limits (below) levaraging ability to refer other elds from your
deployment manifest

```yaml
env:
- name: GOMAXPROCS
  valueFrom:
  resourceFieldRef:
  resource: requests.cpu
```

GOMAXPROCS controls the number of operating system threads that can execute Go code simultaneously. While Go usually manages this well automatically, there are cases where adjustment helps:
- When running in containerized environments with CPU limits (note that Go will detect the number of cores on the underlying node instead of the container's CPU limit)
- When dealing with heavy I/O operations
- When CPU quotas don't align with available cores

- go test -cpuprofile cpu.pprof will run your tests and write a CPU profile to a file named cpu.pprof.
- pprof.StartCPUProfile(w) captures a CPU profile to w that covers the time span until pprof.StopCPUProfile() is called.
- import _ "net/http/pprof" allows you to request a 30s CPU profile by hitting the GET /debug/pprof/profile?seconds=30 endpoint of the default http server that you can start via http.ListenAndServe("localhost:6060", nil)
- runtime.SetCPUProfileRate() lets you to control the sampling rate of the CPU profiler. See CPU Profiler Limitations for current limitations.
- runtime.SetCgoTraceback() can be used to get stack traces into cgo code. benesch/cgosymbolizer has an implementation for Linux and macOS.

file, _ := os.Create("./cpu.pprof")
pprof.StartCPUProfile(file)
defer pprof.StopCPUProfile()

Check how much of overall CPU time is used before going deeper into specific functions which seem to take huge % of CPU time


A few examples of metrics that we can use: `go_sched_latencies`, `go_gc_*`, `go_memstats*`


# Adjusting



# Other tools
3. Load testing - k6 , locust , vegeta
   Especially useful for services with no/minimal downstream dependencies

# Further Reading
- [Go Garbage Collector Guide](https://tip.golang.org/doc/gc-guide) - Comprehensive guide to Go's garbage collector
- Great resource to read more about different types of profiling in Go and their limitations: https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md
- [Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)
- https://github.com/dgryski/go-perfbook


Notes:
- By correlating Go GC pause times with request latency spikes, we identified memory allocation patterns that guided our gRPC middleware optimizations.
- 