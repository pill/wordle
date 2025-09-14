package main

import (
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *DB {
	// Use environment variables for test database
	config := &DatabaseConfig{
		Host:            getEnvString("TEST_DB_HOST", "localhost"),
		Port:            getEnvInt("TEST_DB_PORT", 5432),
		Name:            getEnvString("TEST_DB_NAME", "wordle_test"),
		User:            getEnvString("TEST_DB_USER", "wordle_user"),
		Password:        getEnvString("TEST_DB_PASSWORD", "wordle_password"),
		SSLMode:         "disable",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 15,
	}

	db, err := NewDB(config)
	if err != nil {
		t.Skipf("Skipping database tests: %v", err)
	}

	// Verify required tables exist
	err = db.Migrate()
	if err != nil {
		t.Skipf("Skipping database tests: required tables not found: %v", err)
	}

	return db
}

func TestDatabaseConnection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := db.HealthCheck()
	if err != nil {
		t.Fatalf("Database health check failed: %v", err)
	}
}

func TestGameRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewGameRepository(db)

	// Test CreateGame
	game, err := repo.CreateGame("HELLO", 6)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	if game.ID == "" {
		t.Error("Game ID should not be empty")
	}
	if game.TargetWord != "HELLO" {
		t.Errorf("Expected target word 'HELLO', got '%s'", game.TargetWord)
	}
	if game.MaxGuesses != 6 {
		t.Errorf("Expected max guesses 6, got %d", game.MaxGuesses)
	}
	if game.IsCompleted {
		t.Error("New game should not be completed")
	}
	if game.IsWon {
		t.Error("New game should not be won")
	}
	if game.GuessCount != 0 {
		t.Errorf("New game should have 0 guesses, got %d", game.GuessCount)
	}

	// Test GetGame
	retrievedGame, err := repo.GetGame(game.ID)
	if err != nil {
		t.Fatalf("Failed to get game: %v", err)
	}

	if retrievedGame.ID != game.ID {
		t.Errorf("Expected game ID '%s', got '%s'", game.ID, retrievedGame.ID)
	}
	if retrievedGame.TargetWord != game.TargetWord {
		t.Errorf("Expected target word '%s', got '%s'", game.TargetWord, retrievedGame.TargetWord)
	}

	// Test UpdateGame
	now := time.Now()
	game.IsCompleted = true
	game.IsWon = true
	game.CompletedAt = &now
	game.GuessCount = 3

	err = repo.UpdateGame(game)
	if err != nil {
		t.Fatalf("Failed to update game: %v", err)
	}

	updatedGame, err := repo.GetGame(game.ID)
	if err != nil {
		t.Fatalf("Failed to get updated game: %v", err)
	}

	if !updatedGame.IsCompleted {
		t.Error("Game should be completed after update")
	}
	if !updatedGame.IsWon {
		t.Error("Game should be won after update")
	}
	if updatedGame.GuessCount != 3 {
		t.Errorf("Expected guess count 3, got %d", updatedGame.GuessCount)
	}

	// Test DeleteGame
	err = repo.DeleteGame(game.ID)
	if err != nil {
		t.Fatalf("Failed to delete game: %v", err)
	}

	// Verify game is deleted
	_, err = repo.GetGame(game.ID)
	if err == nil {
		t.Error("Expected error when getting deleted game")
	}
}

func TestGuessRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	gameRepo := NewGameRepository(db)
	guessRepo := NewGuessRepository(db)

	// Create a test game first
	game, err := gameRepo.CreateGame("WORLD", 6)
	if err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	defer gameRepo.DeleteGame(game.ID)

	// Test CreateGuess
	result := GuessResult{
		{Letter: "H", Status: "absent"},
		{Letter: "E", Status: "absent"},
		{Letter: "L", Status: "present"},
		{Letter: "L", Status: "present"},
		{Letter: "O", Status: "correct"},
	}

	guess, err := guessRepo.CreateGuess(game.ID, "HELLO", 1, result)
	if err != nil {
		t.Fatalf("Failed to create guess: %v", err)
	}

	if guess.ID == "" {
		t.Error("Guess ID should not be empty")
	}
	if guess.GameID != game.ID {
		t.Errorf("Expected game ID '%s', got '%s'", game.ID, guess.GameID)
	}
	if guess.GuessWord != "HELLO" {
		t.Errorf("Expected guess word 'HELLO', got '%s'", guess.GuessWord)
	}
	if guess.GuessNumber != 1 {
		t.Errorf("Expected guess number 1, got %d", guess.GuessNumber)
	}
	if len(guess.Result) != 5 {
		t.Errorf("Expected result length 5, got %d", len(guess.Result))
	}

	// Test GetGuess
	retrievedGuess, err := guessRepo.GetGuess(guess.ID)
	if err != nil {
		t.Fatalf("Failed to get guess: %v", err)
	}

	if retrievedGuess.ID != guess.ID {
		t.Errorf("Expected guess ID '%s', got '%s'", guess.ID, retrievedGuess.ID)
	}

	// Test GetGuessesByGameID
	guesses, err := guessRepo.GetGuessesByGameID(game.ID)
	if err != nil {
		t.Fatalf("Failed to get guesses by game ID: %v", err)
	}

	if len(guesses) != 1 {
		t.Errorf("Expected 1 guess, got %d", len(guesses))
	}

	// Test creating multiple guesses
	result2 := GuessResult{
		{Letter: "W", Status: "correct"},
		{Letter: "O", Status: "correct"},
		{Letter: "R", Status: "correct"},
		{Letter: "L", Status: "correct"},
		{Letter: "D", Status: "correct"},
	}

	guess2, err := guessRepo.CreateGuess(game.ID, "WORLD", 2, result2)
	if err != nil {
		t.Fatalf("Failed to create second guess: %v", err)
	}

	// Test GetLatestGuess
	latestGuess, err := guessRepo.GetLatestGuess(game.ID)
	if err != nil {
		t.Fatalf("Failed to get latest guess: %v", err)
	}

	if latestGuess.ID != guess2.ID {
		t.Errorf("Expected latest guess ID '%s', got '%s'", guess2.ID, latestGuess.ID)
	}

	// Test getting all guesses (should be in order)
	allGuesses, err := guessRepo.GetGuessesByGameID(game.ID)
	if err != nil {
		t.Fatalf("Failed to get all guesses: %v", err)
	}

	if len(allGuesses) != 2 {
		t.Errorf("Expected 2 guesses, got %d", len(allGuesses))
	}

	if allGuesses[0].GuessNumber != 1 || allGuesses[1].GuessNumber != 2 {
		t.Error("Guesses should be ordered by guess number")
	}

	// Test DeleteGuess
	err = guessRepo.DeleteGuess(guess.ID)
	if err != nil {
		t.Fatalf("Failed to delete guess: %v", err)
	}

	// Verify guess is deleted
	_, err = guessRepo.GetGuess(guess.ID)
	if err == nil {
		t.Error("Expected error when getting deleted guess")
	}
}

func TestEvaluateGuess(t *testing.T) {
	tests := []struct {
		guess    string
		target   string
		expected []string // status for each letter
	}{
		{
			guess:    "HELLO",
			target:   "HELLO",
			expected: []string{"correct", "correct", "correct", "correct", "correct"},
		},
		{
			guess:    "WORLD",
			target:   "HELLO",
			expected: []string{"absent", "absent", "absent", "present", "present"},
		},
		{
			guess:    "LLAMA",
			target:   "HELLO",
			expected: []string{"absent", "correct", "present", "absent", "absent"},
		},
	}

	for _, test := range tests {
		result := EvaluateGuess(test.guess, test.target)

		if len(result) != len(test.expected) {
			t.Errorf("Expected result length %d, got %d", len(test.expected), len(result))
			continue
		}

		for i, letter := range result {
			if letter.Status != test.expected[i] {
				t.Errorf("For guess '%s' vs target '%s', position %d: expected '%s', got '%s'",
					test.guess, test.target, i, test.expected[i], letter.Status)
			}
		}
	}
}

func TestGameWithGuessesIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	gameRepo := NewGameRepository(db)

	// Create a game
	game, err := gameRepo.CreateGame("CRANE", 6)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer gameRepo.DeleteGame(game.ID)

	guessRepo := NewGuessRepository(db)

	// Add some guesses
	result1 := EvaluateGuess("HELLO", "CRANE")
	_, err = guessRepo.CreateGuess(game.ID, "HELLO", 1, result1)
	if err != nil {
		t.Fatalf("Failed to create first guess: %v", err)
	}

	result2 := EvaluateGuess("CRANE", "CRANE")
	_, err = guessRepo.CreateGuess(game.ID, "CRANE", 2, result2)
	if err != nil {
		t.Fatalf("Failed to create second guess: %v", err)
	}

	// Test GetGameWithGuesses
	gameWithGuesses, err := gameRepo.GetGameWithGuesses(game.ID)
	if err != nil {
		t.Fatalf("Failed to get game with guesses: %v", err)
	}

	if gameWithGuesses.Game.ID != game.ID {
		t.Errorf("Expected game ID '%s', got '%s'", game.ID, gameWithGuesses.Game.ID)
	}

	if len(gameWithGuesses.Guesses) != 2 {
		t.Errorf("Expected 2 guesses, got %d", len(gameWithGuesses.Guesses))
	}

	// Verify guesses are in order
	if gameWithGuesses.Guesses[0].GuessNumber != 1 {
		t.Error("First guess should have guess number 1")
	}
	if gameWithGuesses.Guesses[1].GuessNumber != 2 {
		t.Error("Second guess should have guess number 2")
	}
}
