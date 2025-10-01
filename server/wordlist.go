package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
This could be in the database but for now it's loaded from a file
*/


// WordList represents a collection of words loaded from files
type WordList struct {
	validWords     []string            // All valid words for validation
	validWordSet   map[string]bool     // Set for fast validation lookup
	targetWords    []string            // Common words for game targets
	targetWordSet  map[string]bool     // Set for target word lookup
	validFilePath  string              // Path to validation words file
	targetFilePath string              // Path to target words file
}

// NewWordList creates a new WordList instance
// If validFilePath is empty, it defaults to "valid-wordle-words.txt" in the same directory
// If targetFilePath is empty, it defaults to "common-target-words.txt" in the same directory
func NewWordList(validFilePath string) (*WordList, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Set default paths
	if validFilePath == "" {
		if filepath.Base(dir) == "server" {
			validFilePath = filepath.Join(dir, "valid-wordle-words.txt")
		} else {
			validFilePath = filepath.Join(dir, "server", "valid-wordle-words.txt")
		}
	}

	targetFilePath := ""
	if filepath.Base(dir) == "server" {
		targetFilePath = filepath.Join(dir, "common-target-words.txt")
	} else {
		targetFilePath = filepath.Join(dir, "server", "common-target-words.txt")
	}

	wl := &WordList{
		validFilePath:  validFilePath,
		targetFilePath: targetFilePath,
		validWordSet:   make(map[string]bool),
		targetWordSet:  make(map[string]bool),
	}

	if err := wl.loadWords(); err != nil {
		return nil, err
	}

	return wl, nil
}

// loadWords reads words from both files and populates the word lists
func (wl *WordList) loadWords() error {
	// Load validation words
	if err := wl.loadValidWords(); err != nil {
		return err
	}

	// Load target words
	if err := wl.loadTargetWords(); err != nil {
		return err
	}

	return nil
}

// loadValidWords reads validation words from the file
func (wl *WordList) loadValidWords() error {
	file, err := os.Open(wl.validFilePath)
	if err != nil {
		return fmt.Errorf("failed to open validation word file %s: %w", wl.validFilePath, err)
	}
	defer file.Close()

	wl.validWords = wl.validWords[:0] // Clear existing words
	wl.validWordSet = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wordLower := strings.ToLower(word)
			wl.validWords = append(wl.validWords, wordLower)
			wl.validWordSet[wordLower] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading validation word file: %w", err)
	}

	return nil
}

// loadTargetWords reads target words from the file
func (wl *WordList) loadTargetWords() error {
	file, err := os.Open(wl.targetFilePath)
	if err != nil {
		return fmt.Errorf("failed to open target word file %s: %w", wl.targetFilePath, err)
	}
	defer file.Close()

	wl.targetWords = wl.targetWords[:0] // Clear existing words
	wl.targetWordSet = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wordLower := strings.ToLower(word)
			wl.targetWords = append(wl.targetWords, wordLower)
			wl.targetWordSet[wordLower] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading target word file: %w", err)
	}

	return nil
}

// Size returns the total number of validation words in the list
func (wl *WordList) Size() int {
	return len(wl.validWords)
}

// TargetWordsSize returns the total number of target words in the list
func (wl *WordList) TargetWordsSize() int {
	return len(wl.targetWords)
}

// Contains checks if a word is in the validation list (case-insensitive)
func (wl *WordList) Contains(word string) bool {
	return wl.validWordSet[strings.ToLower(word)]
}

// RandomWord returns a random word from the target words list (for game targets)
func (wl *WordList) RandomWord() string {
	if len(wl.targetWords) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return wl.targetWords[rand.Intn(len(wl.targetWords))]
}

// RandomValidWord returns a random word from the validation list
func (wl *WordList) RandomValidWord() string {
	if len(wl.validWords) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return wl.validWords[rand.Intn(len(wl.validWords))]
}

// WordsOfLength returns all validation words of the specified length
func (wl *WordList) WordsOfLength(length int) []string {
	var result []string
	for _, word := range wl.validWords {
		if len(word) == length {
			result = append(result, word)
		}
	}
	return result
}

// TargetWordsOfLength returns all target words of the specified length
func (wl *WordList) TargetWordsOfLength(length int) []string {
	var result []string
	for _, word := range wl.targetWords {
		if len(word) == length {
			result = append(result, word)
		}
	}
	return result
}

// FiveLetterWords returns all five-letter validation words
func (wl *WordList) FiveLetterWords() []string {
	return wl.WordsOfLength(5)
}

// FiveLetterTargetWords returns all five-letter target words
func (wl *WordList) FiveLetterTargetWords() []string {
	return wl.TargetWordsOfLength(5)
}

// Reload reloads the word list from the file
func (wl *WordList) Reload() error {
	return wl.loadWords()
}

// ToSlice returns a copy of the validation words as a slice
func (wl *WordList) ToSlice() []string {
	result := make([]string, len(wl.validWords))
	copy(result, wl.validWords)
	return result
}

// TargetWordsToSlice returns a copy of the target words as a slice
func (wl *WordList) TargetWordsToSlice() []string {
	result := make([]string, len(wl.targetWords))
	copy(result, wl.targetWords)
	return result
}

// ToSet returns the validation words as a map (set-like structure)
func (wl *WordList) ToSet() map[string]bool {
	result := make(map[string]bool)
	for word := range wl.validWordSet {
		result[word] = true
	}
	return result
}

// TargetWordsToSet returns the target words as a map (set-like structure)
func (wl *WordList) TargetWordsToSet() map[string]bool {
	result := make(map[string]bool)
	for word := range wl.targetWordSet {
		result[word] = true
	}
	return result
}
