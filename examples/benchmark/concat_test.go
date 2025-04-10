package main

import "testing"

var input = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

func BenchmarkConcatenateStrings(b *testing.B) {
	// The loop runs b.N times. b.N is adjusted by the testing framework
	// until the benchmark runs for a stable, measurable duration.
	for i := 0; i < b.N; i++ {
		ConcatenateStrings(input)
	}
}

func BenchmarkConcatenateStringsSlowly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatenateStringsSlowly(input)
	}
}
