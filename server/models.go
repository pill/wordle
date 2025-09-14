package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Game represents a Wordle game session
type Game struct {
	ID          string    `json:"id" db:"id"`
	TargetWord  string    `json:"target_word" db:"target_word"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	IsWon       bool      `json:"is_won" db:"is_won"`
	GuessCount  int       `json:"guess_count" db:"guess_count"`
	MaxGuesses  int       `json:"max_guesses" db:"max_guesses"`
}

// Guess represents a single guess in a game
type Guess struct {
	ID          string      `json:"id" db:"id"`
	GameID      string      `json:"game_id" db:"game_id"`
	GuessWord   string      `json:"guess_word" db:"guess_word"`
	GuessNumber int         `json:"guess_number" db:"guess_number"`
	Result      GuessResult `json:"result" db:"result"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// LetterResult represents the result for a single letter in a guess
type LetterResult struct {
	Letter string `json:"letter"`
	Status string `json:"status"` // "correct", "present", "absent"
}

// GuessResult represents the result of a guess (array of letter results)
type GuessResult []LetterResult

// Value implements the driver.Valuer interface for database storage
func (gr GuessResult) Value() (driver.Value, error) {
	return json.Marshal(gr)
}

// Scan implements the sql.Scanner interface for database retrieval
func (gr *GuessResult) Scan(value interface{}) error {
	if value == nil {
		*gr = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan GuessResult from non-string/[]byte")
	}

	return json.Unmarshal(bytes, gr)
}

// Player represents a player with statistics
type Player struct {
	ID            string    `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	Email         string    `json:"email" db:"email"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	GamesPlayed   int       `json:"games_played" db:"games_played"`
	GamesWon      int       `json:"games_won" db:"games_won"`
	CurrentStreak int       `json:"current_streak" db:"current_streak"`
	MaxStreak     int       `json:"max_streak" db:"max_streak"`
}

// GameStats represents statistics for a game
type GameStats struct {
	ID               string     `json:"id" db:"id"`
	GameID           string     `json:"game_id" db:"game_id"`
	PlayerID         *string    `json:"player_id,omitempty" db:"player_id"`
	WordDifficulty   *float64   `json:"word_difficulty,omitempty" db:"word_difficulty"`
	SolveTimeSeconds *int       `json:"solve_time_seconds,omitempty" db:"solve_time_seconds"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// GameWithGuesses represents a game with all its guesses
type GameWithGuesses struct {
	Game    Game    `json:"game"`
	Guesses []Guess `json:"guesses"`
}

// IsGameComplete checks if the game is complete based on guess count or win status
func (g *Game) IsGameComplete() bool {
	return g.IsWon || g.GuessCount >= g.MaxGuesses
}

// WinRate calculates the win rate for a player
func (p *Player) WinRate() float64 {
	if p.GamesPlayed == 0 {
		return 0.0
	}
	return float64(p.GamesWon) / float64(p.GamesPlayed) * 100
}

// EvaluateGuess evaluates a guess against the target word and returns the result
func EvaluateGuess(guess, target string) GuessResult {
	if len(guess) != len(target) {
		return nil
	}

	guess = strings.ToUpper(guess)
	target = strings.ToUpper(target)

	result := make(GuessResult, len(guess))
	targetChars := make([]rune, len(target))
	copy(targetChars, []rune(target))

	// First pass: mark correct letters
	for i, char := range guess {
		result[i] = LetterResult{
			Letter: string(char),
			Status: "absent",
		}

		if i < len(targetChars) && char == targetChars[i] {
			result[i].Status = "correct"
			targetChars[i] = 0 // Mark as used
		}
	}

	// Second pass: mark present letters
	for i, char := range guess {
		if result[i].Status == "correct" {
			continue
		}

		for j, targetChar := range targetChars {
			if targetChar != 0 && char == targetChar {
				result[i].Status = "present"
				targetChars[j] = 0 // Mark as used
				break
			}
		}
	}

	return result
}

// CreateGameRequest represents a request to create a new game
type CreateGameRequest struct {
	MaxGuesses int `json:"max_guesses,omitempty"`
}

// MakeGuessRequest represents a request to make a guess
type MakeGuessRequest struct {
	GuessWord string `json:"guess_word"`
}

// GameResponse represents a response containing game state
type GameResponse struct {
	Game    Game    `json:"game"`
	Guesses []Guess `json:"guesses,omitempty"`
	Message string  `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
