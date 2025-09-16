package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWordListEdgeCases(t *testing.T) {
	// Create a temporary test file with edge case content
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "edge-case-words.txt")
	
	content := "  HELLO  \n\nWORLD\n  \n\n  CRANE  \n\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	expectedWords := []string{"hello", "world", "crane"}
	if wordList.Size() != len(expectedWords) {
		t.Errorf("Expected %d words, got %d", len(expectedWords), wordList.Size())
	}

	// Verify all expected words are present and normalized
	for _, word := range expectedWords {
		if !wordList.Contains(word) {
			t.Errorf("Expected to find word '%s'", word)
		}
		if !wordList.Contains(strings.ToUpper(word)) {
			t.Errorf("Expected case-insensitive search to find '%s'", strings.ToUpper(word))
		}
	}
}

func TestWordListMixedCase(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "mixed-case-words.txt")
	
	content := "Hello\nWORLD\ncRaNe\nSlAtE\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	// All words should be stored in lowercase
	words := wordList.ToSlice()
	for _, word := range words {
		if word != strings.ToLower(word) {
			t.Errorf("Word '%s' should be stored in lowercase", word)
		}
	}

	// Case-insensitive lookup should work
	testCases := []string{"hello", "HELLO", "Hello", "hELLo"}
	for _, testCase := range testCases {
		if !wordList.Contains(testCase) {
			t.Errorf("Should find '%s' (case insensitive)", testCase)
		}
	}
}

func TestWordListRandomWordDistribution(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "random-test-words.txt")
	
	words := []string{"APPLE", "BANANA", "CHERRY", "DATE", "ELDERBERRY"}
	content := strings.Join(words, "\n")
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	// Test that RandomWord returns valid words
	seenWords := make(map[string]bool)
	for i := 0; i < 20; i++ {
		randomWord := wordList.RandomWord()
		if randomWord == "" {
			t.Error("RandomWord should not return empty string")
		}
		
		if !wordList.Contains(randomWord) {
			t.Errorf("RandomWord returned invalid word: '%s'", randomWord)
		}
		
		seenWords[randomWord] = true
	}

	// We should see some variety (though this is probabilistic)
	if len(seenWords) < 2 {
		t.Logf("Warning: Low randomness detected, only saw %d unique words out of %d", len(seenWords), len(words))
	}
}

func TestWordListWordsOfLengthExtensive(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "length-test-words.txt")
	
	content := "A\nHI\nCAT\nDOG\nHELLO\nWORLD\nCRANE\nSUPERCALIFRAGILISTIC\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	lengthTests := []struct {
		length   int
		expected []string
	}{
		{1, []string{"a"}},
		{2, []string{"hi"}},
		{3, []string{"cat", "dog"}},
		{5, []string{"hello", "world", "crane"}},
		{20, []string{"supercalifragilistic"}},
		{99, []string{}}, // No words of this length
	}

	for _, test := range lengthTests {
		words := wordList.WordsOfLength(test.length)
		if len(words) != len(test.expected) {
			t.Errorf("Length %d: expected %d words, got %d", test.length, len(test.expected), len(words))
			continue
		}

		// Convert to map for easier comparison
		wordMap := make(map[string]bool)
		for _, word := range words {
			wordMap[word] = true
		}

		for _, expectedWord := range test.expected {
			if !wordMap[expectedWord] {
				t.Errorf("Length %d: expected to find word '%s'", test.length, expectedWord)
			}
		}
	}
}

func TestWordListFiveLetterWordsConsistency(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "five-letter-test.txt")
	
	content := "CAT\nHELLO\nWORLD\nCRANE\nSLATE\nAUDIO\nHI\nSUPERLONG\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	fiveLetterWords1 := wordList.FiveLetterWords()
	fiveLetterWords2 := wordList.WordsOfLength(5)

	// These should return the same results
	if len(fiveLetterWords1) != len(fiveLetterWords2) {
		t.Errorf("FiveLetterWords and WordsOfLength(5) returned different counts: %d vs %d", 
			len(fiveLetterWords1), len(fiveLetterWords2))
	}

	// Convert to maps for comparison
	map1 := make(map[string]bool)
	for _, word := range fiveLetterWords1 {
		map1[word] = true
	}

	map2 := make(map[string]bool)
	for _, word := range fiveLetterWords2 {
		map2[word] = true
	}

	for word := range map1 {
		if !map2[word] {
			t.Errorf("Word '%s' in FiveLetterWords but not in WordsOfLength(5)", word)
		}
	}

	for word := range map2 {
		if !map1[word] {
			t.Errorf("Word '%s' in WordsOfLength(5) but not in FiveLetterWords", word)
		}
	}
}

func TestWordListToSetConsistency(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "set-test-words.txt")
	
	content := "HELLO\nWORLD\nCRANE\nHELLO\n" // Duplicate HELLO
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	slice := wordList.ToSlice()
	set := wordList.ToSet()

	// Set should have the same number of unique words
	// Note: slice may contain duplicates, set will not
	expectedUniqueWords := 3 // HELLO, WORLD, CRANE (HELLO appears twice in file)
	if len(set) != expectedUniqueWords {
		t.Errorf("Set size %d should be %d unique words", len(set), expectedUniqueWords)
	}
	
	// The slice should contain all words including duplicates
	expectedTotalWords := 4 // HELLO, WORLD, CRANE, HELLO
	if len(slice) != expectedTotalWords {
		t.Errorf("Slice size %d should be %d total words", len(slice), expectedTotalWords)
	}

	// Every word in slice should be in set
	for _, word := range slice {
		if !set[word] {
			t.Errorf("Word '%s' from slice not found in set", word)
		}
	}

	// Every word in set should be findable via Contains
	for word := range set {
		if !wordList.Contains(word) {
			t.Errorf("Word '%s' from set not found via Contains", word)
		}
	}
}

func TestWordListReloadFunctionality(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "reload-test-words.txt")
	
	// Initial content
	initialContent := "HELLO\nWORLD\n"
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	if wordList.Size() != 2 {
		t.Errorf("Expected 2 words initially, got %d", wordList.Size())
	}

	// Update file content
	updatedContent := "HELLO\nWORLD\nCRANE\nSLATE\n"
	err = os.WriteFile(testFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Size should still be 2 before reload
	if wordList.Size() != 2 {
		t.Errorf("Expected 2 words before reload, got %d", wordList.Size())
	}

	// Reload should pick up new content
	err = wordList.Reload()
	if err != nil {
		t.Fatalf("Failed to reload WordList: %v", err)
	}

	if wordList.Size() != 4 {
		t.Errorf("Expected 4 words after reload, got %d", wordList.Size())
	}

	// New words should be available
	if !wordList.Contains("CRANE") {
		t.Error("Expected to find 'CRANE' after reload")
	}
	if !wordList.Contains("SLATE") {
		t.Error("Expected to find 'SLATE' after reload")
	}
}

func TestWordListFilePathResolution(t *testing.T) {
	// Test with absolute path
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "absolute-path-test.txt")
	
	content := "HELLO\nWORLD\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList with absolute path: %v", err)
	}

	if wordList.Size() != 2 {
		t.Errorf("Expected 2 words with absolute path, got %d", wordList.Size())
	}
}

func TestWordListLargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large-file-test.txt")
	
	// Create a larger file with many words
	var words []string
	for i := 0; i < 1000; i++ {
		// Generate predictable 5-letter words
		word := ""
		num := i
		for j := 0; j < 5; j++ {
			word += string(rune('A' + (num % 26)))
			num /= 26
		}
		words = append(words, word)
	}
	
	content := strings.Join(words, "\n")
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList from large file: %v", err)
	}

	if wordList.Size() != 1000 {
		t.Errorf("Expected 1000 words, got %d", wordList.Size())
	}

	// Test some random lookups
	if !wordList.Contains("AAAAA") {
		t.Error("Expected to find 'AAAAA'")
	}

	// Test performance of RandomWord with large dataset
	for i := 0; i < 100; i++ {
		word := wordList.RandomWord()
		if word == "" {
			t.Error("RandomWord should not return empty string")
		}
	}
}

func TestWordListSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "special-chars-test.txt")
	
	// Include some words with special characters (should be handled gracefully)
	content := "HELLO\nWORLD\nTEST-WORD\nWORD'S\nNORMAL\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	// Should load all words, including those with special characters
	if wordList.Size() != 5 {
		t.Errorf("Expected 5 words, got %d", wordList.Size())
	}

	// Normal words should still work
	if !wordList.Contains("HELLO") {
		t.Error("Expected to find 'HELLO'")
	}
	if !wordList.Contains("NORMAL") {
		t.Error("Expected to find 'NORMAL'")
	}
}
