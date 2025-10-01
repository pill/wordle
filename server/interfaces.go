package main

// Interfaces for dependency injection and testing

// GameRepositoryInterface defines the interface for game repository operations
type GameRepositoryInterface interface {
	CreateGame(targetWord string, maxGuesses int) (*Game, error)
	GetGame(gameID string) (*Game, error)
	UpdateGame(game *Game) error
	DeleteGame(gameID string) error
	GetGameWithGuesses(gameID string) (*GameWithGuesses, error)
	GetRecentGames(limit int) ([]Game, error)
}

// GuessRepositoryInterface defines the interface for guess repository operations
type GuessRepositoryInterface interface {
	CreateGuess(gameID, guessWord string, guessNumber int, result GuessResult) (*Guess, error)
	GetGuess(guessID string) (*Guess, error)
	GetGuessesByGameID(gameID string) ([]Guess, error)
	DeleteGuess(guessID string) error
	GetLatestGuess(gameID string) (*Guess, error)
}

// WordListInterface defines the interface for word list operations
type WordListInterface interface {
	Contains(word string) bool
	RandomWord() string
	RandomValidWord() string
	FiveLetterWords() []string
	FiveLetterTargetWords() []string
	Size() int
	TargetWordsSize() int
}
