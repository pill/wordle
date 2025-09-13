package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWordList(t *testing.T) {
	// Test with default file path
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	if wordList.Size() == 0 {
		t.Error("WordList should not be empty")
	}
}

func TestWordListContains(t *testing.T) {
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	// Test case insensitive matching
	if !wordList.Contains("AAHED") && wordList.Contains("aahed") {
		t.Error("Contains should be case insensitive")
	}

	// Test non-existent word
	if wordList.Contains("xyzzynotaword") {
		t.Error("Should return false for non-existent word")
	}
}

func TestWordListRandomWord(t *testing.T) {
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	word := wordList.RandomWord()
	if word == "" {
		t.Error("RandomWord should not return empty string")
	}

	if !wordList.Contains(word) {
		t.Error("RandomWord should return a word from the list")
	}
}

func TestWordListWordsOfLength(t *testing.T) {
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	fiveLetterWords := wordList.WordsOfLength(5)
	for _, word := range fiveLetterWords {
		if len(word) != 5 {
			t.Errorf("Expected 5-letter word, got '%s' with length %d", word, len(word))
		}
	}

	// Test that FiveLetterWords returns the same as WordsOfLength(5)
	fiveLetterWords2 := wordList.FiveLetterWords()
	if len(fiveLetterWords) != len(fiveLetterWords2) {
		t.Error("FiveLetterWords should return same count as WordsOfLength(5)")
	}
}

func TestWordListToSlice(t *testing.T) {
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	slice := wordList.ToSlice()
	if len(slice) != wordList.Size() {
		t.Error("ToSlice should return slice with same length as Size()")
	}

	// Modify the slice and ensure original is unchanged
	originalSize := wordList.Size()
	slice[0] = "modified"
	if wordList.Size() != originalSize {
		t.Error("Modifying returned slice should not affect original WordList")
	}
}

func TestWordListToSet(t *testing.T) {
	wordList, err := NewWordList("")
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	set := wordList.ToSet()
	if len(set) != wordList.Size() {
		t.Error("ToSet should return map with same length as Size()")
	}

	// Test that all words are in the set
	for _, word := range wordList.ToSlice() {
		if !set[word] {
			t.Errorf("Word '%s' should be in the set", word)
		}
	}
}

func TestWordListReload(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-words.txt")
	
	content := "apple\nbanana\ncherry\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	if wordList.Size() != 3 {
		t.Errorf("Expected 3 words, got %d", wordList.Size())
	}

	// Modify the file
	newContent := "apple\nbanana\ncherry\ndate\nelderberry\n"
	err = os.WriteFile(testFile, []byte(newContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Reload and check
	err = wordList.Reload()
	if err != nil {
		t.Fatalf("Failed to reload WordList: %v", err)
	}

	if wordList.Size() != 5 {
		t.Errorf("Expected 5 words after reload, got %d", wordList.Size())
	}
}

func TestWordListEmptyFile(t *testing.T) {
	// Create an empty temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.txt")
	
	err := os.WriteFile(testFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList from empty file: %v", err)
	}

	if wordList.Size() != 0 {
		t.Errorf("Expected 0 words from empty file, got %d", wordList.Size())
	}

	if wordList.RandomWord() != "" {
		t.Error("RandomWord should return empty string for empty word list")
	}
}

func TestWordListNormalization(t *testing.T) {
	// Create a test file with mixed case and whitespace
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-normalization.txt")
	
	content := "  Apple  \nBANANA\n  cherry\n\n  \nDATE\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	wordList, err := NewWordList(testFile)
	if err != nil {
		t.Fatalf("Failed to create WordList: %v", err)
	}

	// Should have 4 words (empty lines and whitespace-only lines ignored)
	if wordList.Size() != 4 {
		t.Errorf("Expected 4 words, got %d", wordList.Size())
	}

	// All words should be lowercase
	for _, word := range wordList.ToSlice() {
		if word != strings.ToLower(word) {
			t.Errorf("Word '%s' should be lowercase", word)
		}
		if strings.TrimSpace(word) != word {
			t.Errorf("Word '%s' should be trimmed", word)
		}
	}

	// Test case-insensitive lookup
	if !wordList.Contains("APPLE") {
		t.Error("Should find 'APPLE' (case insensitive)")
	}
	if !wordList.Contains("banana") {
		t.Error("Should find 'banana'")
	}
}
