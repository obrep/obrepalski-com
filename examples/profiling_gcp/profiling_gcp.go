package main

import (
	"log"
	"os"

	"cloud.google.com/go/profiler"
)

func main() {
	// Configuration for the profiler.
	cfg := profiler.Config{
		Service:        "your-service-name",      // Replace with your service name
		ServiceVersion: os.Getenv("APP_VERSION"), // Use an env var for version (e.g., BUILD_ID, git SHA)
		// ProjectID is optional if running on GCP infra (inferred)
		// ProjectID: "your-gcp-project-id",
	}

	// Start the profiler. Errors are logged if it fails to start.
	if err := profiler.Start(cfg); err != nil {
		log.Fatalf("WARN: Failed to start profiler: %v", err)
		// Usually, you wouldn't stop the app if the profiler fails,
		// so we just log the error.
	}

	// ... rest of your application startup and logic ...
	log.Println("Application started...")
}
