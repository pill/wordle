# WordList Go Package

A simple Go package for loading and managing the Wordle word list.

## Usage

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Create a new word list instance
    wordList, err := NewWordList("")
    if err != nil {
        log.Fatal(err)
    }

    // Get basic information
    fmt.Printf("Total words: %d\n", wordList.Size())
    fmt.Printf("Random word: %s\n", wordList.RandomWord())

    // Check if a word is valid
    if wordList.Contains("hello") {
        fmt.Println("'hello' is a valid word")
    }

    // Get words of specific length
    fiveLetterWords := wordList.FiveLetterWords()
    fmt.Printf("Number of five-letter words: %d\n", len(fiveLetterWords))

    // Get words of any length
    threeLetterWords := wordList.WordsOfLength(3)
    fmt.Printf("Number of three-letter words: %d\n", len(threeLetterWords))
}
```

## Methods

### Constructor
- `NewWordList(filePath string) (*WordList, error)` - Creates a new WordList instance
  - If filePath is empty, defaults to "valid-wordle-words.txt" in the current directory structure

### Core Methods
- `Size() int` - Returns the total number of words
- `Contains(word string) bool` - Checks if a word is in the list (case-insensitive)
- `RandomWord() string` - Returns a random word from the list
- `WordsOfLength(length int) []string` - Returns all words of the specified length
- `FiveLetterWords() []string` - Returns all five-letter words
- `Reload() error` - Reloads the word list from the file

### Utility Methods
- `ToSlice() []string` - Returns the words as a slice (copy)
- `ToSet() map[string]bool` - Returns the words as a map for set-like operations

## Running the Demo

To see the WordList in action, run:

```bash
go run main.go wordlist.go
```

This will show:
- Total word count
- Random word samples
- Five-letter word count and samples
- Word validation examples
- Word length distribution
- Sample words by length

## Running Tests

To run the comprehensive test suite:

```bash
go test -v
```

This runs tests for:
- Word loading and validation
- Case-insensitive matching
- Random word generation
- Length filtering
- File reloading
- Error handling
- Data normalization

## Building

To build the package:

```bash
go build
```

To build and run:

```bash
go run main.go wordlist.go
```

## File Structure

- `valid-wordle-words.txt` - The source word list file (14,855 five-letter words)
- `wordlist.go` - The WordList struct implementation
- `main.go` - Demo program showing usage
- `wordlist_test.go` - Comprehensive test suite
- `README.md` - This documentation file

## Features

- **Case-insensitive**: All word operations are case-insensitive
- **Memory efficient**: Uses both slice and map for different access patterns
- **Thread-safe reads**: Once loaded, the word list can be safely read from multiple goroutines
- **Comprehensive testing**: Full test coverage including edge cases
- **Flexible file paths**: Automatically detects whether running from root or server directory
- **Error handling**: Proper error handling for file operations
- **Data normalization**: Automatically trims whitespace and converts to lowercase

## Performance

- Loading 14,855 words: ~0.1 seconds
- Word lookup (Contains): O(1) average case
- Random word selection: O(1)
- Length filtering: O(n) where n is total words
- Memory usage: ~1MB for the full word list