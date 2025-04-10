package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	// Import pprof package for side effects: registers HTTP handlers.
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"
)

const (
	// Number of iterations for the CPU-heavy calculation
	piIterations = 50000
	// Number of report lines to generate (causes allocations)
	reportLines = 1000
	// How often to process a "request"
	processInterval = 50 * time.Millisecond
)

func main() {
	// Start the pprof HTTP server on localhost:6060
	// Run this in a separate goroutine so it doesn't block.
	go func() {
		log.Println("Starting pprof server on localhost:6060")
		// Normally you might want to trap the error, but for a demo just log
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Starting simulated work...")

	// Simulate processing requests periodically
	ticker := time.NewTicker(processInterval)
	defer ticker.Stop()

	for range ticker.C {
		processRequest(rand.Intn(1000)) // Simulate passing some request ID
	}
}

// processRequest simulates handling a top-level request.
// It calls functions to perform different kinds of work.
func processRequest(requestID int) {
	log.Printf("Processing request %d\n", requestID)

	// Simulate some validation step that involves heavy computation
	isValid := validateRequest(requestID)
	if !isValid {
		log.Printf("Request %d validation failed (simulated)\n", requestID)
		return // Don't proceed if validation fails
	}

	// Simulate generating some data or report which involves allocations
	report := generateReportData(requestID)

	// (Imagine doing something with the report here)
	_ = report // Use report to avoid compiler error

	log.Printf("Finished processing request %d\n", requestID)
}

// validateRequest simulates a validation process.
// It calls a CPU-intensive function.
func validateRequest(id int) bool {
	// Simulate some quick checks first
	if id < 0 {
		return false
	}
	// Call the function that consumes CPU
	result := calculatePiApprox(piIterations)

	// Use the result to make it seem important
	// (In a real scenario, the result might influence validation)
	// Here, we just ensure it's > 0, which it always will be.
	return result > 0
}

// calculatePiApprox is a deliberately CPU-intensive function.
// It uses the Leibniz formula for Pi (inefficiently) as an example.
func calculatePiApprox(iterations int) float64 {
	var pi float64
	for i := 0; i < iterations; i++ {
		pi += math.Pow(-1, float64(i)) / (2*float64(i) + 1)
	}
	return pi * 4 // Leibniz formula converges to pi/4
}

// generateReportData simulates creating some data structure or report.
// It calls a function that causes memory allocations through string operations.
func generateReportData(id int) string {
	var report strings.Builder // Use strings.Builder for efficiency *usually*,
	// but the helper function will be inefficient for demo purposes.

	report.WriteString("Report for Request ID: ")
	report.WriteString(strconv.Itoa(id))
	report.WriteString("\n")
	report.WriteString("---------------------------\n")

	// Generate multiple lines, each causing allocations in the helper
	for i := 0; i < reportLines; i++ {
		report.WriteString(generateLineItem(i))
		report.WriteString("\n")
	}

	report.WriteString("---------------------------\n")
	report.WriteString("End of Report\n")

	return report.String()
}

// generateLineItem simulates creating a single line item for the report.
// This function uses inefficient string concatenation (+) which causes allocations.
func generateLineItem(itemNumber int) string {
	line := "Item: " + strconv.Itoa(itemNumber) + " | "
	line += "Value: " + fmt.Sprintf("%.2f", rand.Float64()*100) + " | "
	line += "Status: OK"
	return line
}