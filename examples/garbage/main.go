package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Poorly aligned struct - fields are not ordered by size
// This causes memory padding and wastes space
type User struct {
	active      bool   // 1 byte + 7 bytes padding
	id          int64  // 8 bytes
	isVIP       bool   // 1 byte + 3 bytes padding
	score       int32  // 4 bytes
	enabled     bool   // 1 byte + 7 bytes padding
	lastLoginMs int64  // 8 bytes
	name        string // 16 bytes (string header)
	tags        []string // 24 bytes (slice header)
	metadata    map[string]interface{} // 8 bytes (pointer)
}

// Another poorly aligned struct
type Transaction struct {
	processed   bool   // 1 byte + 7 bytes padding
	amount      float64 // 8 bytes
	isPending   bool   // 1 byte + 3 bytes padding
	retryCount  int32  // 4 bytes
	completed   bool   // 1 byte + 7 bytes padding
	timestamp   int64  // 8 bytes
	id          string // 16 bytes
	description string // 16 bytes
	metadata    interface{} // 16 bytes (interface)
}

type DataProcessor struct {
	users         []User
	transactions  []Transaction
	cache         map[string]interface{} // Using interface{} causes boxing
	mu            sync.Mutex
	tempData      [][]byte // Will accumulate garbage
	sessionData   map[int]map[string]string // Nested maps are inefficient
}

var (
	processor *DataProcessor
	stopChan  = make(chan bool)
)

func main() {
	processor = &DataProcessor{
		users:        make([]User, 0),
		transactions: make([]Transaction, 0),
		cache:        make(map[string]interface{}),
		sessionData:  make(map[int]map[string]string),
	}

	// Start pprof server
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println("CPU profile: http://localhost:6060/debug/pprof/profile?seconds=30")
		log.Println("Heap profile: http://localhost:6060/debug/pprof/heap")
		log.Println("Allocations: http://localhost:6060/debug/pprof/allocs")
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	// Print initial memory stats
	printMemStats("Initial")

	// Start multiple worker goroutines to create GC pressure
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}

	// Start a goroutine that creates memory leaks
	go leakyGoroutine()

	// Start periodic stats printer
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			printMemStats("Periodic")
		}
	}()

	// Main processing loop
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	requestID := 0
	for {
		select {
		case <-ticker.C:
			requestID++
			processRequest(requestID)
		case <-stopChan:
			return
		}
	}
}

func worker(id int) {
	for {
		select {
		case <-stopChan:
			return
		default:
			// Inefficient string concatenation in a loop
			result := ""
			for i := 0; i < 100; i++ {
				result = result + "worker" + strconv.Itoa(id) + "_item" + strconv.Itoa(i) + "_"
			}
			
			// Create temporary data that will be garbage collected
			tempData := make([]byte, 1024*10) // 10KB allocation
			for i := range tempData {
				tempData[i] = byte(rand.Intn(256))
			}
			
			// Use the data to prevent optimization
			if len(result) > 0 && len(tempData) > 0 {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func processRequest(id int) {
	// Create user with inefficient patterns
	user := createUser(id)
	processor.mu.Lock()
	processor.users = append(processor.users, user)
	
	// Keep only last 1000 users (inefficient removal from slice beginning)
	if len(processor.users) > 1000 {
		processor.users = processor.users[1:] // Creates a new backing array
	}
	processor.mu.Unlock()

	// Process transaction with allocations
	transaction := createTransaction(id)
	processTransaction(transaction)

	// Inefficient cache update
	updateCache(id)

	// Generate report with many allocations
	report := generateReport(id)
	_ = report // Use to avoid optimization
}

func createUser(id int) User {
	// Inefficient: creating string arrays with concatenation
	tags := make([]string, 0)
	for i := 0; i < 20; i++ {
		tag := "tag_" + strconv.Itoa(id) + "_" + strconv.Itoa(i)
		tags = append(tags, tag) // Growing slice causes reallocations
	}

	// Creating metadata with boxing
	metadata := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		key := "key_" + strconv.Itoa(i)
		if i%2 == 0 {
			metadata[key] = i // Boxing int to interface{}
		} else {
			metadata[key] = "value_" + strconv.Itoa(i) // Boxing string
		}
	}

	return User{
		id:          int64(id),
		name:        "User_" + strconv.Itoa(id), // String concatenation
		active:      id%2 == 0,
		isVIP:       id%10 == 0,
		enabled:     true,
		score:       int32(rand.Intn(1000)),
		lastLoginMs: time.Now().UnixMilli(),
		tags:        tags,
		metadata:    metadata,
	}
}

func createTransaction(id int) Transaction {
	// Inefficient string building
	desc := ""
	for i := 0; i < 10; i++ {
		desc = desc + "Part" + strconv.Itoa(i) + " "
	}

	return Transaction{
		id:          "TX_" + strconv.Itoa(id),
		amount:      rand.Float64() * 1000,
		processed:   false,
		isPending:   true,
		completed:   false,
		retryCount:  0,
		timestamp:   time.Now().Unix(),
		description: desc,
		metadata:    map[string]string{"type": "standard", "category": "online"},
	}
}

func processTransaction(tx Transaction) {
	processor.mu.Lock()
	defer processor.mu.Unlock()

	// Inefficient: append and potential reallocation
	processor.transactions = append(processor.transactions, tx)

	// Inefficient: JSON marshaling creates allocations
	data, _ := json.Marshal(tx)
	
	// Store in temp data (accumulates garbage)
	processor.tempData = append(processor.tempData, data)
	
	// Keep only last 500 transactions (inefficient)
	if len(processor.transactions) > 500 {
		// Creating new slice instead of reusing
		processor.transactions = processor.transactions[len(processor.transactions)-500:]
	}
	
	// Clean temp data inefficiently
	if len(processor.tempData) > 100 {
		processor.tempData = make([][]byte, 0) // Discards old slice entirely
	}
}

func updateCache(id int) {
	processor.mu.Lock()
	defer processor.mu.Unlock()

	// Inefficient: string concatenation for keys
	for i := 0; i < 50; i++ {
		key := "cache_" + strconv.Itoa(id) + "_" + strconv.Itoa(i)
		
		// Boxing different types
		if i%3 == 0 {
			processor.cache[key] = i
		} else if i%3 == 1 {
			processor.cache[key] = "value_" + strconv.Itoa(i)
		} else {
			processor.cache[key] = map[string]int{"nested": i}
		}
	}

	// Inefficient cache cleanup
	if len(processor.cache) > 1000 {
		// Creating new map instead of deleting entries
		newCache := make(map[string]interface{})
		count := 0
		for k, v := range processor.cache {
			if count < 500 {
				newCache[k] = v
				count++
			}
		}
		processor.cache = newCache
	}
}

func generateReport(id int) string {
	// Very inefficient string concatenation
	report := "Report for request: " + strconv.Itoa(id) + "\n"
	report = report + "========================\n"
	
	// Get current users (creates a copy)
	processor.mu.Lock()
	usersCopy := make([]User, len(processor.users))
	copy(usersCopy, processor.users)
	processor.mu.Unlock()
	
	// Build report with string concatenation
	for _, user := range usersCopy {
		report = report + "User ID: " + strconv.FormatInt(user.id, 10) + "\n"
		report = report + "Name: " + user.name + "\n"
		report = report + "Score: " + strconv.Itoa(int(user.score)) + "\n"
		
		// Concatenate all tags
		tagsStr := ""
		for _, tag := range user.tags {
			tagsStr = tagsStr + tag + ", "
		}
		report = report + "Tags: " + tagsStr + "\n"
		report = report + "---\n"
	}
	
	// Add transaction summary
	processor.mu.Lock()
	txCount := len(processor.transactions)
	processor.mu.Unlock()
	
	report = report + "Total transactions: " + strconv.Itoa(txCount) + "\n"
	report = report + "========================\n"
	
	return report
}

func leakyGoroutine() {
	// This simulates a memory leak pattern
	leakedData := make(map[int][]byte)
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	id := 0
	for range ticker.C {
		id++
		// Allocate memory that's never freed
		data := make([]byte, 1024*100) // 100KB
		for i := range data {
			data[i] = byte(rand.Intn(256))
		}
		leakedData[id] = data
		
		// Simulate forgetting to clean old entries
		// This map will grow indefinitely
		
		// Also create session data that's never cleaned
		processor.mu.Lock()
		processor.sessionData[id] = make(map[string]string)
		for i := 0; i < 10; i++ {
			key := "session_" + strconv.Itoa(i)
			value := "data_" + strconv.Itoa(id) + "_" + strconv.Itoa(i)
			processor.sessionData[id][key] = value
		}
		processor.mu.Unlock()
	}
}

func printMemStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	log.Printf("[%s] Memory Stats:\n", label)
	log.Printf("  Alloc = %v MB", m.Alloc/1024/1024)
	log.Printf("  TotalAlloc = %v MB", m.TotalAlloc/1024/1024)
	log.Printf("  Sys = %v MB", m.Sys/1024/1024)
	log.Printf("  NumGC = %v", m.NumGC)
	log.Printf("  HeapAlloc = %v MB", m.HeapAlloc/1024/1024)
	log.Printf("  HeapInuse = %v MB", m.HeapInuse/1024/1024)
	log.Printf("  HeapObjects = %v", m.HeapObjects)
	
	// Force a GC to see the impact
	runtime.GC()
	
	runtime.ReadMemStats(&m)
	log.Printf("  After GC - HeapAlloc = %v MB", m.HeapAlloc/1024/1024)
	log.Printf("  After GC - NumGC = %v", m.NumGC)
	fmt.Println()
}