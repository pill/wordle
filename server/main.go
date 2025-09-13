package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// Create a new word list instance
	wordList, err := NewWordList("")
	if err != nil {
		fmt.Printf("Error creating word list: %v\n", err)
		return
	}

	fmt.Println("=== WordList Test ===")
	fmt.Printf("Total words loaded: %d\n", wordList.Size())
	
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	fmt.Printf("Random word: %s\n", wordList.RandomWord())
	fmt.Printf("Random word: %s\n", wordList.RandomWord())
	fmt.Printf("Random word: %s\n", wordList.RandomWord())

	fmt.Println("\n=== Five Letter Words ===")
	fiveLetterWords := wordList.FiveLetterWords()
	fmt.Printf("Number of five-letter words: %d\n", len(fiveLetterWords))
	
	// Get sample five-letter words
	sampleWords := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(fiveLetterWords); i++ {
		idx := rand.Intn(len(fiveLetterWords))
		sampleWords = append(sampleWords, fiveLetterWords[idx])
	}
	fmt.Printf("Sample five-letter words: %s\n", strings.Join(sampleWords, ", "))

	fmt.Println("\n=== Word Validation ===")
	testWords := []string{"hello", "world", "apple", "xyzzy", "abask", "aahed"}
	for _, word := range testWords {
		valid := wordList.Contains(word)
		status := "not valid"
		if valid {
			status = "valid"
		}
		fmt.Printf("'%s' is %s\n", word, status)
	}

	fmt.Println("\n=== Word Length Distribution ===")
	for length := 3; length <= 8; length++ {
		count := len(wordList.WordsOfLength(length))
		fmt.Printf("%d-letter words: %d\n", length, count)
	}

	fmt.Println("\n=== Sample Words by Length ===")
	for length := 3; length <= 8; length++ {
		words := wordList.WordsOfLength(length)
		if len(words) > 0 {
			sampleSize := 5
			if len(words) < sampleSize {
				sampleSize = len(words)
			}
			
			sample := make([]string, 0, sampleSize)
			for i := 0; i < sampleSize; i++ {
				idx := rand.Intn(len(words))
				sample = append(sample, words[idx])
			}
			fmt.Printf("%d-letter words sample: %s\n", length, strings.Join(sample, ", "))
		}
	}
}
