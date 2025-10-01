package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// Mock implementations for testing

type MockGameRepository struct {
	games         map[string]*Game
	nextID        int
	shouldFailGet bool
	shouldFailSave bool
}

func NewMockGameRepository() *MockGameRepository {
	return &MockGameRepository{
		games:  make(map[string]*Game),
		nextID: 1,
	}
}

func (m *MockGameRepository) CreateGame(targetWord string, maxGuesses int) (*Game, error) {
	if m.shouldFailSave {
		return nil, errors.New("mock save error")
	}

	id := string(rune(m.nextID + 64)) // Convert to letter (A, B, C, etc.)
	m.nextID++

	game := &Game{
		ID:          id,
		TargetWord:  targetWord,
		CreatedAt:   time.Now(),
		IsCompleted: false,
		IsWon:       false,
		GuessCount:  0,
		MaxGuesses:  maxGuesses,
	}

	m.games[id] = game
	return game, nil
}

func (m *MockGameRepository) GetGame(gameID string) (*Game, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get error")
	}

	game, exists := m.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Return a copy to avoid modification issues
	gameCopy := *game
	return &gameCopy, nil
}

func (m *MockGameRepository) UpdateGame(game *Game) error {
	if m.shouldFailSave {
		return errors.New("mock update error")
	}

	_, exists := m.games[game.ID]
	if !exists {
		return errors.New("game not found")
	}

	// Store a copy
	gameCopy := *game
	m.games[game.ID] = &gameCopy
	return nil
}

func (m *MockGameRepository) GetGameWithGuesses(gameID string) (*GameWithGuesses, error) {
	game, err := m.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	// For testing, return empty guesses
	return &GameWithGuesses{
		Game:    *game,
		Guesses: []Guess{},
	}, nil
}

func (m *MockGameRepository) DeleteGame(gameID string) error {
	if m.shouldFailSave {
		return errors.New("mock delete error")
	}

	_, exists := m.games[gameID]
	if !exists {
		return errors.New("game not found")
	}

	delete(m.games, gameID)
	return nil
}

func (m *MockGameRepository) GetRecentGames(limit int) ([]Game, error) {
	var games []Game
	for _, game := range m.games {
		games = append(games, *game)
		if len(games) >= limit {
			break
		}
	}
	return games, nil
}

type MockGuessRepository struct {
	guesses         map[string][]Guess
	shouldFailSave  bool
	shouldFailGet   bool
	nextGuessID     int
}

func NewMockGuessRepository() *MockGuessRepository {
	return &MockGuessRepository{
		guesses:     make(map[string][]Guess),
		nextGuessID: 1,
	}
}

func (m *MockGuessRepository) CreateGuess(gameID, guessWord string, guessNumber int, result GuessResult) (*Guess, error) {
	if m.shouldFailSave {
		return nil, errors.New("mock save guess error")
	}

	// Check for duplicate guess numbers
	if guesses, exists := m.guesses[gameID]; exists {
		for _, guess := range guesses {
			if guess.GuessNumber == guessNumber {
				return nil, errors.New("guess number already exists")
			}
		}
	}

	guess := &Guess{
		ID:          string(rune(m.nextGuessID + 64)),
		GameID:      gameID,
		GuessWord:   guessWord,
		GuessNumber: guessNumber,
		Result:      result,
		CreatedAt:   time.Now(),
	}
	m.nextGuessID++

	if m.guesses[gameID] == nil {
		m.guesses[gameID] = []Guess{}
	}
	m.guesses[gameID] = append(m.guesses[gameID], *guess)

	return guess, nil
}

func (m *MockGuessRepository) GetGuessesByGameID(gameID string) ([]Guess, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get guesses error")
	}

	guesses, exists := m.guesses[gameID]
	if !exists {
		return []Guess{}, nil
	}

	// Sort by guess number
	sortedGuesses := make([]Guess, len(guesses))
	copy(sortedGuesses, guesses)
	
	// Simple bubble sort for testing
	for i := 0; i < len(sortedGuesses)-1; i++ {
		for j := 0; j < len(sortedGuesses)-i-1; j++ {
			if sortedGuesses[j].GuessNumber > sortedGuesses[j+1].GuessNumber {
				sortedGuesses[j], sortedGuesses[j+1] = sortedGuesses[j+1], sortedGuesses[j]
			}
		}
	}

	return sortedGuesses, nil
}

func (m *MockGuessRepository) GetGuess(guessID string) (*Guess, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get guess error")
	}

	// Simple implementation for testing
	for _, guesses := range m.guesses {
		for _, guess := range guesses {
			if guess.ID == guessID {
				guessCopy := guess
				return &guessCopy, nil
			}
		}
	}
	return nil, errors.New("guess not found")
}

func (m *MockGuessRepository) DeleteGuess(guessID string) error {
	if m.shouldFailSave {
		return errors.New("mock delete guess error")
	}

	// Find and remove the guess
	for gameID, guesses := range m.guesses {
		for i, guess := range guesses {
			if guess.ID == guessID {
				// Remove from slice
				m.guesses[gameID] = append(guesses[:i], guesses[i+1:]...)
				return nil
			}
		}
	}
	return errors.New("guess not found")
}

func (m *MockGuessRepository) GetLatestGuess(gameID string) (*Guess, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get latest guess error")
	}

	guesses, exists := m.guesses[gameID]
	if !exists || len(guesses) == 0 {
		return nil, errors.New("no guesses found")
	}

	// Find the guess with the highest guess number
	var latest *Guess
	for _, guess := range guesses {
		if latest == nil || guess.GuessNumber > latest.GuessNumber {
			guessCopy := guess
			latest = &guessCopy
		}
	}

	return latest, nil
}

type MockWordList struct {
	words         []string
	shouldFailGet bool
}

func NewMockWordList() *MockWordList {
	return &MockWordList{
		words: []string{"HELLO", "WORLD", "CRANE", "SLATE", "AUDIO", "QUICK", "BROWN"},
	}
}

func (m *MockWordList) Contains(word string) bool {
	if m.shouldFailGet {
		return false
	}
	
	word = strings.ToUpper(word)
	for _, w := range m.words {
		if w == word {
			return true
		}
	}
	return false
}

func (m *MockWordList) RandomWord() string {
	if len(m.words) == 0 {
		return ""
	}
	return m.words[0] // Always return first word for predictable testing
}

func (m *MockWordList) FiveLetterWords() []string {
	return m.words
}

func (m *MockWordList) Size() int {
	return len(m.words)
}

func (m *MockWordList) RandomValidWord() string {
	if len(m.words) == 0 {
		return ""
	}
	return m.words[1] // Return second word for testing
}

func (m *MockWordList) FiveLetterTargetWords() []string {
	return m.words // For testing, use same words as target words
}

func (m *MockWordList) TargetWordsSize() int {
	return len(m.words)
}

// Test functions

func TestGameServiceCreateNewGame(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("CreateNewGame should not return error: %v", err)
	}

	if game.ID == "" {
		t.Error("Game should have an ID")
	}
	if game.TargetWord != "HELLO" { // First word from mock
		t.Errorf("Expected target word 'HELLO', got '%s'", game.TargetWord)
	}
	if game.MaxGuesses != 6 {
		t.Errorf("Expected max guesses 6, got %d", game.MaxGuesses)
	}
	if game.GuessCount != 0 {
		t.Errorf("New game should have 0 guesses, got %d", game.GuessCount)
	}
	if game.IsCompleted {
		t.Error("New game should not be completed")
	}
	if game.IsWon {
		t.Error("New game should not be won")
	}
}

func TestGameServiceCreateNewGameNoWords(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := &MockWordList{words: []string{}} // Empty word list
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	_, err := service.CreateNewGame()
	if err == nil {
		t.Error("Expected error when no words available")
	}
	if !strings.Contains(err.Error(), "no five-letter words available") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGameServiceMakeGuessValid(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create a game first
	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Make a valid guess
	response, err := service.MakeGuess(game.ID, "WORLD")
	if err != nil {
		t.Fatalf("MakeGuess should not return error: %v", err)
	}

	if response.Game.GuessCount != 1 {
		t.Errorf("Expected guess count 1, got %d", response.Game.GuessCount)
	}
	if response.Game.IsWon {
		t.Error("Game should not be won with incorrect guess")
	}
	if response.Game.IsCompleted {
		t.Error("Game should not be completed after one guess")
	}
	if !strings.Contains(response.Message, "5 guess(es) remaining") {
		t.Errorf("Expected remaining guesses message, got: %s", response.Message)
	}
}

func TestGameServiceMakeGuessWinning(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create a game
	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Make winning guess (same as target word)
	response, err := service.MakeGuess(game.ID, "HELLO")
	if err != nil {
		t.Fatalf("MakeGuess should not return error: %v", err)
	}

	if response.Game.GuessCount != 1 {
		t.Errorf("Expected guess count 1, got %d", response.Game.GuessCount)
	}
	if !response.Game.IsWon {
		t.Error("Game should be won with correct guess")
	}
	if !response.Game.IsCompleted {
		t.Error("Game should be completed when won")
	}
	if !strings.Contains(response.Message, "Congratulations") {
		t.Errorf("Expected congratulations message, got: %s", response.Message)
	}
}

func TestGameServiceMakeGuessInvalidWord(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create a game
	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Try invalid word
	_, err = service.MakeGuess(game.ID, "ZZZZZ")
	if err == nil {
		t.Error("Expected error for invalid word")
	}
	if !strings.Contains(err.Error(), "not a valid word") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGameServiceMakeGuessWrongLength(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create a game
	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Try wrong length word
	_, err = service.MakeGuess(game.ID, "HI")
	if err == nil {
		t.Error("Expected error for wrong length word")
	}
	if !strings.Contains(err.Error(), "must be 5 letters long") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGameServiceMakeGuessGameNotFound(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Try to make guess on non-existent game
	_, err := service.MakeGuess("nonexistent", "HELLO")
	if err == nil {
		t.Error("Expected error for non-existent game")
	}
	if !strings.Contains(err.Error(), "failed to get game") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGameServiceMakeGuessGameCompleted(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create and complete a game
	game, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Manually mark game as completed
	game.IsCompleted = true
	gameRepo.UpdateGame(game)

	// Try to make guess on completed game
	_, err = service.MakeGuess(game.ID, "WORLD")
	if err == nil {
		t.Error("Expected error for completed game")
	}
	if !strings.Contains(err.Error(), "already completed") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGameServiceValidateWord(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Test valid word
	if !service.ValidateWord("HELLO") {
		t.Error("Expected 'HELLO' to be valid")
	}

	// Test invalid word
	if service.ValidateWord("ZZZZZ") {
		t.Error("Expected 'ZZZZZ' to be invalid")
	}

	// Test wrong length
	if service.ValidateWord("HI") {
		t.Error("Expected 'HI' to be invalid due to length")
	}

	// Test with whitespace
	if !service.ValidateWord(" HELLO ") {
		t.Error("Expected 'HELLO' with whitespace to be valid")
	}
}

func TestGameServiceGetGameStats(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	stats, err := service.GetGameStats()
	if err != nil {
		t.Fatalf("GetGameStats should not return error: %v", err)
	}

	expectedStats := map[string]interface{}{
		"total_words":        7, // From mock word list
		"five_letter_words":  7,
		"max_guesses":       6,
		"word_length":       5,
	}

	for key, expected := range expectedStats {
		if stats[key] != expected {
			t.Errorf("Expected %s to be %v, got %v", key, expected, stats[key])
		}
	}
}

func TestGameServiceGetRecentGames(t *testing.T) {
	gameRepo := NewMockGameRepository()
	guessRepo := NewMockGuessRepository()
	wordList := NewMockWordList()
	config := &GameConfig{MaxGuesses: 6, WordLength: 5}

	service := NewGameServiceWithInterfaces(gameRepo, guessRepo, wordList, config)

	// Create some games
	_, err := service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create first game: %v", err)
	}
	_, err = service.CreateNewGame()
	if err != nil {
		t.Fatalf("Failed to create second game: %v", err)
	}

	// Test with valid limit
	games, err := service.GetRecentGames(10)
	if err != nil {
		t.Fatalf("GetRecentGames should not return error: %v", err)
	}

	if len(games) != 2 {
		t.Errorf("Expected 2 games, got %d", len(games))
	}

	// Test with limit bounds
	games, err = service.GetRecentGames(0)
	if err != nil {
		t.Fatalf("GetRecentGames should not return error: %v", err)
	}
	// Should default to 10
	if len(games) > 10 {
		t.Errorf("Expected at most 10 games with limit 0, got %d", len(games))
	}

	games, err = service.GetRecentGames(200)
	if err != nil {
		t.Fatalf("GetRecentGames should not return error: %v", err)
	}
	// Should default to 10
	if len(games) > 10 {
		t.Errorf("Expected at most 10 games with limit 200, got %d", len(games))
	}
}
