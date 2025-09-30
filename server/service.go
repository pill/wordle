package main

import (
	"fmt"
	"strings"
	"time"
)

// GameService handles business logic for Wordle games
type GameService struct {
	gameRepo  GameRepositoryInterface
	guessRepo GuessRepositoryInterface
	wordList  WordListInterface
	config    *GameConfig
}

// NewGameService creates a new game service
func NewGameService(db *DB, wordList *WordList, config *GameConfig) *GameService {
	return &GameService{
		gameRepo:  NewGameRepository(db),
		guessRepo: NewGuessRepository(db),
		wordList:  wordList,
		config:    config,
	}
}

// NewGameServiceWithInterfaces creates a new game service with injectable interfaces
func NewGameServiceWithInterfaces(gameRepo GameRepositoryInterface, guessRepo GuessRepositoryInterface, wordList WordListInterface, config *GameConfig) *GameService {
	return &GameService{
		gameRepo:  gameRepo,
		guessRepo: guessRepo,
		wordList:  wordList,
		config:    config,
	}
}

// CreateNewGame creates a new game with a random target word
func (s *GameService) CreateNewGame() (*Game, error) {
	// Get a random five-letter word
	// TODO: this could be in the database but for now it's loaded from a file
	// TODO: random word should not repeat for user
	fiveLetterWords := s.wordList.FiveLetterWords()
	if len(fiveLetterWords) == 0 {
		return nil, fmt.Errorf("no five-letter words available")
	}

	targetWord := strings.ToUpper(s.wordList.RandomWord())
	maxGuesses := s.config.MaxGuesses

	game, err := s.gameRepo.CreateGame(targetWord, maxGuesses)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

// GetGame retrieves a game by ID
func (s *GameService) GetGame(gameID string) (*Game, error) {
	return s.gameRepo.GetGame(gameID)
}

// GetGameWithGuesses retrieves a game with all its guesses
func (s *GameService) GetGameWithGuesses(gameID string) (*GameWithGuesses, error) {
	return s.gameRepo.GetGameWithGuesses(gameID)
}

// MakeGuess processes a guess for a game
func (s *GameService) MakeGuess(gameID, guessWord string) (*GameResponse, error) {
	// Get the current game
	game, err := s.gameRepo.GetGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Check if game is already completed
	if game.IsCompleted {
		return nil, fmt.Errorf("game is already completed")
	}

	// Validate guess word
	guessWord = strings.ToUpper(strings.TrimSpace(guessWord))
	if len(guessWord) != s.config.WordLength {
		return nil, fmt.Errorf("guess must be %d letters long", s.config.WordLength)
	}

	// Check if word is valid
	if !s.wordList.Contains(guessWord) {
		return nil, fmt.Errorf("'%s' is not a valid word", guessWord)
	}

	// Check if player has remaining guesses
	if game.GuessCount >= game.MaxGuesses {
		return nil, fmt.Errorf("no remaining guesses")
	}

	// Evaluate the guess
	result := EvaluateGuess(guessWord, game.TargetWord)
	guessNumber := game.GuessCount + 1

	// Create the guess record
	_, err = s.guessRepo.CreateGuess(gameID, guessWord, guessNumber, result)
	if err != nil {
		return nil, fmt.Errorf("failed to save guess: %w", err)
	}

	// Update game state
	game.GuessCount = guessNumber
	isWin := guessWord == game.TargetWord
	game.IsWon = isWin
	game.IsCompleted = isWin || game.GuessCount >= game.MaxGuesses

	if game.IsCompleted {
		now := time.Now()
		game.CompletedAt = &now
	}

	// Save updated game
	err = s.gameRepo.UpdateGame(game)
	if err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	// Get all guesses for response
	guesses, err := s.guessRepo.GetGuessesByGameID(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guesses: %w", err)
	}

	// Prepare response message
	var message string
	if game.IsWon {
		message = fmt.Sprintf("Congratulations! You won in %d guess(es)!", game.GuessCount)
	} else if game.IsCompleted {
		message = fmt.Sprintf("Game over! The word was '%s'", game.TargetWord)
	} else {
		remaining := game.MaxGuesses - game.GuessCount
		message = fmt.Sprintf("Good guess! %d guess(es) remaining", remaining)
	}

	return &GameResponse{
		Game:    *game,
		Guesses: guesses,
		Message: message,
	}, nil
}

// GetRecentGames gets recent games
func (s *GameService) GetRecentGames(limit int) ([]Game, error) {
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}
	return s.gameRepo.GetRecentGames(limit)
}

// DeleteGame deletes a game
func (s *GameService) DeleteGame(gameID string) error {
	return s.gameRepo.DeleteGame(gameID)
}

// ValidateWord checks if a word is valid for Wordle
func (s *GameService) ValidateWord(word string) bool {
	word = strings.TrimSpace(word)
	if len(word) != s.config.WordLength {
		return false
	}
	return s.wordList.Contains(word)
}

// GetGameStats returns basic statistics about games
func (s *GameService) GetGameStats() (map[string]interface{}, error) {
	// This could be expanded with more sophisticated statistics
	stats := make(map[string]interface{})

	// For now, just return basic word list info
	stats["total_words"] = s.wordList.Size()
	stats["five_letter_words"] = len(s.wordList.FiveLetterWords())
	stats["max_guesses"] = s.config.MaxGuesses
	stats["word_length"] = s.config.WordLength

	return stats, nil
}
