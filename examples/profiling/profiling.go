package main

import (
	"log"
	"net/http"

	// Import pprof package for side effects: registers HTTP handlers.
	// We use the blank identifier _ because we only need the side effects (handler registration)
	// from its init() function, not any functions directly from the package.
	_ "net/http/pprof"
)

func main() {
	// Start the pprof HTTP server on a separate port and goroutine.
	// Running it in a separate goroutine ensures it doesn't block the main application logic.
	// Using a different port (e.g., 6060) is common practice to avoid interfering
	// with the main application's port (e.g., 8080).
	go func() {
		log.Println("Starting pprof server on localhost:6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Your main service logic would go here...
	// For demonstration, we'll just block forever.
	log.Println("Main application running...")
	select {}
}
