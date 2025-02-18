package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// Add debug endpoint for pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Your service code here
}
