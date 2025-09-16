package main

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"
)

func TestGameIsGameComplete(t *testing.T) {
	tests := []struct {
		name        string
		game        Game
		expected    bool
		description string
	}{
		{
			name: "Game won",
			game: Game{
				IsWon:      true,
				GuessCount: 3,
				MaxGuesses: 6,
			},
			expected:    true,
			description: "Game should be complete when won",
		},
		{
			name: "Game lost - max guesses reached",
			game: Game{
				IsWon:      false,
				GuessCount: 6,
				MaxGuesses: 6,
			},
			expected:    true,
			description: "Game should be complete when max guesses reached",
		},
		{
			name: "Game lost - exceeded max guesses",
			game: Game{
				IsWon:      false,
				GuessCount: 7,
				MaxGuesses: 6,
			},
			expected:    true,
			description: "Game should be complete when guesses exceed max",
		},
		{
			name: "Game in progress",
			game: Game{
				IsWon:      false,
				GuessCount: 3,
				MaxGuesses: 6,
			},
			expected:    false,
			description: "Game should not be complete when still in progress",
		},
		{
			name: "New game",
			game: Game{
				IsWon:      false,
				GuessCount: 0,
				MaxGuesses: 6,
			},
			expected:    false,
			description: "New game should not be complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.game.IsGameComplete()
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
			}
		})
	}
}

func TestPlayerWinRate(t *testing.T) {
	tests := []struct {
		name     string
		player   Player
		expected float64
	}{
		{
			name: "Perfect win rate",
			player: Player{
				GamesPlayed: 10,
				GamesWon:    10,
			},
			expected: 100.0,
		},
		{
			name: "Zero win rate",
			player: Player{
				GamesPlayed: 10,
				GamesWon:    0,
			},
			expected: 0.0,
		},
		{
			name: "50% win rate",
			player: Player{
				GamesPlayed: 10,
				GamesWon:    5,
			},
			expected: 50.0,
		},
		{
			name: "No games played",
			player: Player{
				GamesPlayed: 0,
				GamesWon:    0,
			},
			expected: 0.0,
		},
		{
			name: "Partial win rate",
			player: Player{
				GamesPlayed: 3,
				GamesWon:    1,
			},
			expected: 33.333333333333336, // 1/3 * 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.player.WinRate()
			// Use a small tolerance for floating point comparison
			tolerance := 0.000001
			if result < tt.expected-tolerance || result > tt.expected+tolerance {
				t.Errorf("Expected win rate %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestEvaluateGuess(t *testing.T) {
	tests := []struct {
		name     string
		guess    string
		target   string
		expected []LetterResult
	}{
		{
			name:   "Perfect match",
			guess:  "HELLO",
			target: "HELLO",
			expected: []LetterResult{
				{Letter: "H", Status: "correct"},
				{Letter: "E", Status: "correct"},
				{Letter: "L", Status: "correct"},
				{Letter: "L", Status: "correct"},
				{Letter: "O", Status: "correct"},
			},
		},
		{
			name:   "No matches",
			guess:  "ABCDE",
			target: "FGHIJ",
			expected: []LetterResult{
				{Letter: "A", Status: "absent"},
				{Letter: "B", Status: "absent"},
				{Letter: "C", Status: "absent"},
				{Letter: "D", Status: "absent"},
				{Letter: "E", Status: "absent"},
			},
		},
		{
			name:   "Mixed results",
			guess:  "WORLD",
			target: "HELLO",
			expected: []LetterResult{
				{Letter: "W", Status: "absent"},
				{Letter: "O", Status: "present"},
				{Letter: "R", Status: "absent"},
				{Letter: "L", Status: "correct"},
				{Letter: "D", Status: "absent"},
			},
		},
		{
			name:   "Duplicate letters in guess",
			guess:  "LLAMA",
			target: "HELLO",
			expected: []LetterResult{
				{Letter: "L", Status: "present"},
				{Letter: "L", Status: "present"},
				{Letter: "A", Status: "absent"},
				{Letter: "M", Status: "absent"},
				{Letter: "A", Status: "absent"},
			},
		},
		{
			name:   "Case insensitive",
			guess:  "hello",
			target: "HELLO",
			expected: []LetterResult{
				{Letter: "H", Status: "correct"},
				{Letter: "E", Status: "correct"},
				{Letter: "L", Status: "correct"},
				{Letter: "L", Status: "correct"},
				{Letter: "O", Status: "correct"},
			},
		},
		{
			name:   "Complex duplicate handling",
			guess:  "SPEED",
			target: "ERASE",
			expected: []LetterResult{
				{Letter: "S", Status: "present"},
				{Letter: "P", Status: "absent"},
				{Letter: "E", Status: "present"},
				{Letter: "E", Status: "present"},
				{Letter: "D", Status: "absent"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateGuess(tt.guess, tt.target)

			if len(result) != len(tt.expected) {
				t.Fatalf("Expected result length %d, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if result[i].Letter != expected.Letter {
					t.Errorf("Position %d: expected letter '%s', got '%s'", i, expected.Letter, result[i].Letter)
				}
				if result[i].Status != expected.Status {
					t.Errorf("Position %d (%s): expected status '%s', got '%s'", i, expected.Letter, expected.Status, result[i].Status)
				}
			}
		})
	}
}

func TestEvaluateGuessInvalidLength(t *testing.T) {
	// Test with different lengths
	result := EvaluateGuess("HELLO", "HI")
	if result != nil {
		t.Error("Expected nil result for mismatched word lengths")
	}

	result = EvaluateGuess("HI", "HELLO")
	if result != nil {
		t.Error("Expected nil result for mismatched word lengths")
	}
}

func TestGuessResultValue(t *testing.T) {
	result := GuessResult{
		{Letter: "H", Status: "correct"},
		{Letter: "E", Status: "present"},
		{Letter: "L", Status: "absent"},
	}

	value, err := result.Value()
	if err != nil {
		t.Fatalf("Value() should not return error: %v", err)
	}

	// Should be able to unmarshal the JSON
	var unmarshaled GuessResult
	err = json.Unmarshal(value.([]byte), &unmarshaled)
	if err != nil {
		t.Fatalf("Should be able to unmarshal JSON: %v", err)
	}

	if len(unmarshaled) != len(result) {
		t.Errorf("Expected length %d, got %d", len(result), len(unmarshaled))
	}

	for i, expected := range result {
		if unmarshaled[i].Letter != expected.Letter {
			t.Errorf("Position %d: expected letter '%s', got '%s'", i, expected.Letter, unmarshaled[i].Letter)
		}
		if unmarshaled[i].Status != expected.Status {
			t.Errorf("Position %d: expected status '%s', got '%s'", i, expected.Status, unmarshaled[i].Status)
		}
	}
}

func TestGuessResultScan(t *testing.T) {
	// Test scanning from []byte
	jsonData := `[{"letter":"H","status":"correct"},{"letter":"E","status":"present"}]`
	var result GuessResult

	err := result.Scan([]byte(jsonData))
	if err != nil {
		t.Fatalf("Scan from []byte should not return error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	if result[0].Letter != "H" || result[0].Status != "correct" {
		t.Error("First letter result not scanned correctly")
	}
	if result[1].Letter != "E" || result[1].Status != "present" {
		t.Error("Second letter result not scanned correctly")
	}

	// Test scanning from string
	var result2 GuessResult
	err = result2.Scan(jsonData)
	if err != nil {
		t.Fatalf("Scan from string should not return error: %v", err)
	}

	if len(result2) != 2 {
		t.Errorf("Expected length 2, got %d", len(result2))
	}

	// Test scanning from nil
	var result3 GuessResult
	err = result3.Scan(nil)
	if err != nil {
		t.Fatalf("Scan from nil should not return error: %v", err)
	}
	if result3 != nil {
		t.Error("Expected nil result when scanning nil")
	}

	// Test scanning from invalid type
	var result4 GuessResult
	err = result4.Scan(123)
	if err == nil {
		t.Error("Expected error when scanning from invalid type")
	}

	// Test scanning invalid JSON
	var result5 GuessResult
	err = result5.Scan("invalid json")
	if err == nil {
		t.Error("Expected error when scanning invalid JSON")
	}
}

func TestGuessResultDriverValuer(t *testing.T) {
	result := GuessResult{
		{Letter: "H", Status: "correct"},
	}

	// Test that GuessResult implements driver.Valuer
	var _ driver.Valuer = result

	value, err := result.Value()
	if err != nil {
		t.Fatalf("Value() should not return error: %v", err)
	}

	// Should return []byte
	_, ok := value.([]byte)
	if !ok {
		t.Error("Value() should return []byte")
	}
}

func TestCreateGameRequest(t *testing.T) {
	// Test that the struct can be created and marshaled
	request := CreateGameRequest{
		MaxGuesses: 8,
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Should be able to marshal CreateGameRequest: %v", err)
	}

	var unmarshaled CreateGameRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Should be able to unmarshal CreateGameRequest: %v", err)
	}

	if unmarshaled.MaxGuesses != 8 {
		t.Errorf("Expected MaxGuesses 8, got %d", unmarshaled.MaxGuesses)
	}
}

func TestMakeGuessRequest(t *testing.T) {
	request := MakeGuessRequest{
		GuessWord: "HELLO",
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Should be able to marshal MakeGuessRequest: %v", err)
	}

	var unmarshaled MakeGuessRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Should be able to unmarshal MakeGuessRequest: %v", err)
	}

	if unmarshaled.GuessWord != "HELLO" {
		t.Errorf("Expected GuessWord 'HELLO', got '%s'", unmarshaled.GuessWord)
	}
}

func TestGameResponse(t *testing.T) {
	now := time.Now()
	response := GameResponse{
		Game: Game{
			ID:          "test-id",
			TargetWord:  "HELLO",
			CreatedAt:   now,
			IsCompleted: false,
			IsWon:       false,
			GuessCount:  2,
			MaxGuesses:  6,
		},
		Guesses: []Guess{
			{
				ID:          "guess-1",
				GameID:      "test-id",
				GuessWord:   "WORLD",
				GuessNumber: 1,
				Result: GuessResult{
					{Letter: "W", Status: "absent"},
				},
				CreatedAt: now,
			},
		},
		Message: "Good guess!",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Should be able to marshal GameResponse: %v", err)
	}

	var unmarshaled GameResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Should be able to unmarshal GameResponse: %v", err)
	}

	if unmarshaled.Game.ID != "test-id" {
		t.Errorf("Expected Game ID 'test-id', got '%s'", unmarshaled.Game.ID)
	}
	if len(unmarshaled.Guesses) != 1 {
		t.Errorf("Expected 1 guess, got %d", len(unmarshaled.Guesses))
	}
	if unmarshaled.Message != "Good guess!" {
		t.Errorf("Expected message 'Good guess!', got '%s'", unmarshaled.Message)
	}
}

func TestErrorResponse(t *testing.T) {
	response := ErrorResponse{
		Error:   "Test error",
		Code:    400,
		Details: "Test details",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Should be able to marshal ErrorResponse: %v", err)
	}

	var unmarshaled ErrorResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Should be able to unmarshal ErrorResponse: %v", err)
	}

	if unmarshaled.Error != "Test error" {
		t.Errorf("Expected error 'Test error', got '%s'", unmarshaled.Error)
	}
	if unmarshaled.Code != 400 {
		t.Errorf("Expected code 400, got %d", unmarshaled.Code)
	}
	if unmarshaled.Details != "Test details" {
		t.Errorf("Expected details 'Test details', got '%s'", unmarshaled.Details)
	}
}
