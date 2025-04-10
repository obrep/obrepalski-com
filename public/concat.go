package main

import "strings"

func ConcatenateStrings(parts []string) string {
	return strings.Join(parts, "")
}

// Slower implementation for comparison:
func ConcatenateStringsSlowly(parts []string) string {
	var result string
	for _, s := range parts {
		result += s
	}
	return result
}
