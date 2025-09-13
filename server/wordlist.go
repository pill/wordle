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

// WordList represents a collection of words loaded from a file
type WordList struct {
	words    []string
	wordSet  map[string]bool
	filePath string
}

// NewWordList creates a new WordList instance
// If filePath is empty, it defaults to "valid-wordle-words.txt" in the same directory
func NewWordList(filePath string) (*WordList, error) {
	if filePath == "" {
		// Get the directory of the current executable/source
		dir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		// Check if we're in the server directory or the root directory
		if filepath.Base(dir) == "server" {
			filePath = filepath.Join(dir, "valid-wordle-words.txt")
		} else {
			filePath = filepath.Join(dir, "server", "valid-wordle-words.txt")
		}
	}

	wl := &WordList{
		filePath: filePath,
		wordSet:  make(map[string]bool),
	}

	if err := wl.loadWords(); err != nil {
		return nil, err
	}

	return wl, nil
}

// loadWords reads words from the file and populates the word list
func (wl *WordList) loadWords() error {
	file, err := os.Open(wl.filePath)
	if err != nil {
		return fmt.Errorf("failed to open word file %s: %w", wl.filePath, err)
	}
	defer file.Close()

	wl.words = wl.words[:0] // Clear existing words
	wl.wordSet = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wordLower := strings.ToLower(word)
			wl.words = append(wl.words, wordLower)
			wl.wordSet[wordLower] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading word file: %w", err)
	}

	return nil
}

// Size returns the total number of words in the list
func (wl *WordList) Size() int {
	return len(wl.words)
}

// Contains checks if a word is in the list (case-insensitive)
func (wl *WordList) Contains(word string) bool {
	return wl.wordSet[strings.ToLower(word)]
}

// RandomWord returns a random word from the list
func (wl *WordList) RandomWord() string {
	if len(wl.words) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return wl.words[rand.Intn(len(wl.words))]
}

// WordsOfLength returns all words of the specified length
func (wl *WordList) WordsOfLength(length int) []string {
	var result []string
	for _, word := range wl.words {
		if len(word) == length {
			result = append(result, word)
		}
	}
	return result
}

// FiveLetterWords returns all five-letter words
func (wl *WordList) FiveLetterWords() []string {
	return wl.WordsOfLength(5)
}

// Reload reloads the word list from the file
func (wl *WordList) Reload() error {
	return wl.loadWords()
}

// ToSlice returns a copy of the words as a slice
func (wl *WordList) ToSlice() []string {
	result := make([]string, len(wl.words))
	copy(result, wl.words)
	return result
}

// ToSet returns the words as a map (set-like structure)
func (wl *WordList) ToSet() map[string]bool {
	result := make(map[string]bool)
	for word := range wl.wordSet {
		result[word] = true
	}
	return result
}
