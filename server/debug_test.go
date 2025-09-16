package main

import (
	"fmt"
	"testing"
)

func TestDebugEvaluateGuess(t *testing.T) {
	// Test the actual behavior
	result := EvaluateGuess("LLAMA", "HELLO")
	fmt.Println("LLAMA vs HELLO:")
	for i, r := range result {
		fmt.Printf("Position %d: %s -> %s\n", i, r.Letter, r.Status)
	}
	
	// Let's also test the word list issue
	fmt.Println("\nTesting duplicate word handling...")
}
