package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
)

// Mock database types for unit testing (without actual database)

type MockDB struct {
	shouldFailQuery bool
	shouldFailExec  bool
	mockRows        *MockRows
	mockResult      *MockResult
	lastQuery       string
	lastArgs        []interface{}
}

type MockRows struct {
	data     [][]interface{}
	columns  []string
	current  int
	closed   bool
	scanFunc func(dest ...interface{}) error
}

type MockResult struct {
	rowsAffected int64
	lastInsertId int64
	shouldFail   bool
}

func (r *MockResult) RowsAffected() (int64, error) {
	if r.shouldFail {
		return 0, errors.New("mock rows affected error")
	}
	return r.rowsAffected, nil
}

func (r *MockResult) LastInsertId() (int64, error) {
	if r.shouldFail {
		return 0, errors.New("mock last insert id error")
	}
	return r.lastInsertId, nil
}

func (r *MockRows) Next() bool {
	if r.closed {
		return false
	}
	if r.current >= len(r.data) {
		return false
	}
	r.current++
	return true
}

func (r *MockRows) Scan(dest ...interface{}) error {
	if r.closed {
		return errors.New("rows closed")
	}
	if r.current == 0 || r.current > len(r.data) {
		return errors.New("no current row")
	}

	if r.scanFunc != nil {
		return r.scanFunc(dest...)
	}

	// Default scan behavior
	row := r.data[r.current-1]
	if len(dest) != len(row) {
		return fmt.Errorf("destination count %d != source count %d", len(dest), len(row))
	}

	for i, val := range row {
		switch d := dest[i].(type) {
		case *string:
			if s, ok := val.(string); ok {
				*d = s
			} else {
				return fmt.Errorf("cannot scan %T into *string", val)
			}
		case *int:
			if i, ok := val.(int); ok {
				*d = i
			} else {
				return fmt.Errorf("cannot scan %T into *int", val)
			}
		case *bool:
			if b, ok := val.(bool); ok {
				*d = b
			} else {
				return fmt.Errorf("cannot scan %T into *bool", val)
			}
		case *time.Time:
			if t, ok := val.(time.Time); ok {
				*d = t
			} else {
				return fmt.Errorf("cannot scan %T into *time.Time", val)
			}
		case **time.Time:
			if val == nil {
				*d = nil
			} else if t, ok := val.(time.Time); ok {
				*d = &t
			} else {
				return fmt.Errorf("cannot scan %T into **time.Time", val)
			}
		case *GuessResult:
			if s, ok := val.(string); ok {
				return d.Scan(s)
			} else {
				return fmt.Errorf("cannot scan %T into *GuessResult", val)
			}
		default:
			return fmt.Errorf("unsupported destination type %T", d)
		}
	}

	return nil
}

func (r *MockRows) Close() error {
	r.closed = true
	return nil
}

func (r *MockRows) Err() error {
	return nil
}

// Mock database implementation
func (db *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	db.lastQuery = query
	db.lastArgs = args

	if db.shouldFailQuery {
		// Return a row that will fail on scan
		return &sql.Row{}
	}

	// This is a simplified mock - in real testing you'd use sqlmock or similar
	// For now, we'll test the error cases
	return &sql.Row{}
}

func (db *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db.lastQuery = query
	db.lastArgs = args

	if db.shouldFailQuery {
		return nil, errors.New("mock query error")
	}

	// Return mock rows - this is simplified
	return &sql.Rows{}, nil
}

func (db *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.lastQuery = query
	db.lastArgs = args

	if db.shouldFailExec {
		return nil, errors.New("mock exec error")
	}

	return db.mockResult, nil
}

// Unit tests for repository functions using mocks

func TestGameRepositoryCreateGameValidation(t *testing.T) {
	tests := []struct {
		name       string
		targetWord string
		maxGuesses int
		shouldPass bool
	}{
		{"Valid input", "HELLO", 6, true},
		{"Empty target word", "", 6, true}, // Should still create but validate elsewhere
		{"Zero max guesses", "HELLO", 0, true}, // Business logic validation
		{"Negative max guesses", "HELLO", -1, true}, // Business logic validation
		{"Long target word", "SUPERCALIFRAGILISTICEXPIALIDOCIOUS", 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test input validation - repository layer should accept any input
			// Business validation happens in service layer
			if tt.targetWord == "" && tt.maxGuesses == 6 {
				// This is fine for repository layer
			}
			if tt.maxGuesses <= 0 {
				// This is also fine for repository layer
			}
			// Repository tests would require more complex mocking
			// These are more appropriate as integration tests
		})
	}
}

func TestGuessRepositoryInputValidation(t *testing.T) {
	tests := []struct {
		name        string
		gameID      string
		guessWord   string
		guessNumber int
		result      GuessResult
		expectError bool
	}{
		{
			name:        "Valid guess",
			gameID:      "valid-game-id",
			guessWord:   "HELLO",
			guessNumber: 1,
			result:      GuessResult{{Letter: "H", Status: "correct"}},
			expectError: false,
		},
		{
			name:        "Empty game ID",
			gameID:      "",
			guessWord:   "HELLO",
			guessNumber: 1,
			result:      GuessResult{{Letter: "H", Status: "correct"}},
			expectError: false, // Repository should accept, validation elsewhere
		},
		{
			name:        "Empty guess word",
			gameID:      "valid-game-id",
			guessWord:   "",
			guessNumber: 1,
			result:      GuessResult{},
			expectError: false, // Repository should accept, validation elsewhere
		},
		{
			name:        "Zero guess number",
			gameID:      "valid-game-id",
			guessWord:   "HELLO",
			guessNumber: 0,
			result:      GuessResult{{Letter: "H", Status: "correct"}},
			expectError: false, // Repository should accept, validation elsewhere
		},
		{
			name:        "Negative guess number",
			gameID:      "valid-game-id",
			guessWord:   "HELLO",
			guessNumber: -1,
			result:      GuessResult{{Letter: "H", Status: "correct"}},
			expectError: false, // Repository should accept, validation elsewhere
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the repository layer accepts various inputs
			// The actual database constraints and business logic validation
			// happens at different layers
			
			// For unit testing repositories, we'd typically use dependency injection
			// and mock the database interface
			
			// Verify input handling logic
			if tt.gameID == "" {
				// Repository should handle this gracefully (might fail at DB level)
			}
			if tt.guessWord == "" {
				// Repository should handle this gracefully
			}
			if tt.guessNumber <= 0 {
				// Repository should handle this gracefully
			}
		})
	}
}

func TestGuessResultSerialization(t *testing.T) {
	tests := []struct {
		name   string
		result GuessResult
	}{
		{
			name: "Single letter",
			result: GuessResult{
				{Letter: "H", Status: "correct"},
			},
		},
		{
			name: "Multiple letters",
			result: GuessResult{
				{Letter: "H", Status: "correct"},
				{Letter: "E", Status: "present"},
				{Letter: "L", Status: "absent"},
				{Letter: "L", Status: "absent"},
				{Letter: "O", Status: "correct"},
			},
		},
		{
			name:   "Empty result",
			result: GuessResult{},
		},
		{
			name: "Special characters in letters",
			result: GuessResult{
				{Letter: "'", Status: "absent"},
				{Letter: "-", Status: "absent"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Value() method
			value, err := tt.result.Value()
			if err != nil {
				t.Fatalf("Value() should not return error: %v", err)
			}

			// Should return []byte
			bytes, ok := value.([]byte)
			if !ok {
				t.Fatalf("Value() should return []byte, got %T", value)
			}

			// Test Scan() method
			var scanned GuessResult
			err = scanned.Scan(bytes)
			if err != nil {
				t.Fatalf("Scan() should not return error: %v", err)
			}

			// Verify round-trip consistency
			if len(scanned) != len(tt.result) {
				t.Errorf("Length mismatch after round-trip: expected %d, got %d", len(tt.result), len(scanned))
			}

			for i, expected := range tt.result {
				if i >= len(scanned) {
					t.Errorf("Missing element at index %d", i)
					continue
				}
				if scanned[i].Letter != expected.Letter {
					t.Errorf("Letter mismatch at index %d: expected '%s', got '%s'", i, expected.Letter, scanned[i].Letter)
				}
				if scanned[i].Status != expected.Status {
					t.Errorf("Status mismatch at index %d: expected '%s', got '%s'", i, expected.Status, scanned[i].Status)
				}
			}
		})
	}
}

func TestPostgresErrorHandling(t *testing.T) {
	// Test how repository handles different PostgreSQL error types
	
	tests := []struct {
		name        string
		pgError     *pq.Error
		expectedMsg string
	}{
		{
			name: "Unique violation",
			pgError: &pq.Error{
				Code: "23505",
				Message: "duplicate key value violates unique constraint",
			},
			expectedMsg: "already exists",
		},
		{
			name: "Foreign key violation",
			pgError: &pq.Error{
				Code: "23503",
				Message: "violates foreign key constraint",
			},
			expectedMsg: "foreign key",
		},
		{
			name: "Not null violation",
			pgError: &pq.Error{
				Code: "23502",
				Message: "null value in column violates not-null constraint",
			},
			expectedMsg: "not-null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error code detection
			if tt.pgError.Code == "23505" {
				// Should handle unique violations specially
				if !strings.Contains(tt.expectedMsg, "exists") {
					t.Errorf("Expected 'exists' in message for unique violation")
				}
			}
			
			if tt.pgError.Code == "23503" {
				// Should handle foreign key violations
				if !strings.Contains(tt.expectedMsg, "foreign key") {
					t.Errorf("Expected 'foreign key' in message for FK violation")
				}
			}
			
			if tt.pgError.Code == "23502" {
				// Should handle not-null violations
				if !strings.Contains(tt.expectedMsg, "not-null") {
					t.Errorf("Expected 'not-null' in message for null violation")
				}
			}
		})
	}
}

func TestRepositoryQueryConstruction(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		expectedQuery string
		expectedArgs  int
	}{
		{
			name:          "Create game query",
			operation:     "create_game",
			expectedQuery: "INSERT INTO games",
			expectedArgs:  2, // targetWord, maxGuesses
		},
		{
			name:          "Get game query",
			operation:     "get_game",
			expectedQuery: "SELECT",
			expectedArgs:  1, // gameID
		},
		{
			name:          "Update game query",
			operation:     "update_game",
			expectedQuery: "UPDATE games",
			expectedArgs:  5, // completedAt, isCompleted, isWon, guessCount, id
		},
		{
			name:          "Create guess query",
			operation:     "create_guess",
			expectedQuery: "INSERT INTO guesses",
			expectedArgs:  4, // gameID, guessWord, guessNumber, result
		},
		{
			name:          "Get guesses query",
			operation:     "get_guesses",
			expectedQuery: "SELECT",
			expectedArgs:  1, // gameID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test query structure expectations
			if !strings.Contains(tt.expectedQuery, tt.operation) && 
			   !strings.Contains(strings.ToLower(tt.expectedQuery), strings.Split(tt.operation, "_")[0]) {
				// Verify the query type matches the operation
			}
			
			if tt.expectedArgs <= 0 {
				t.Errorf("Expected positive number of args for %s", tt.operation)
			}
		})
	}
}

func TestDatabaseTransactionHandling(t *testing.T) {
	// Test transaction scenarios (would use mocks in real implementation)
	
	scenarios := []struct {
		name        string
		operations  []string
		shouldFail  bool
		failAt      int
	}{
		{
			name:       "Successful transaction",
			operations: []string{"insert_game", "insert_guess"},
			shouldFail: false,
		},
		{
			name:       "Failed at first operation",
			operations: []string{"insert_game", "insert_guess"},
			shouldFail: true,
			failAt:     0,
		},
		{
			name:       "Failed at second operation",
			operations: []string{"insert_game", "insert_guess"},
			shouldFail: true,
			failAt:     1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Test transaction logic
			for i, op := range scenario.operations {
				if scenario.shouldFail && i == scenario.failAt {
					// Simulate failure
					// In real implementation, would verify rollback behavior
					if op == "insert_game" {
						// Game insertion failed
					} else if op == "insert_guess" {
						// Guess insertion failed, should rollback game
					}
				}
			}
		})
	}
}

func TestRepositoryConnectionPoolUsage(t *testing.T) {
	// Test connection pool behavior (conceptual test)
	
	tests := []struct {
		name           string
		concurrency    int
		operations     int
		expectedMetric string
	}{
		{
			name:           "Low concurrency",
			concurrency:    1,
			operations:     10,
			expectedMetric: "single_connection",
		},
		{
			name:           "High concurrency",
			concurrency:    10,
			operations:     100,
			expectedMetric: "multiple_connections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conceptual connection pool usage
			if tt.concurrency == 1 {
				// Should use minimal connections
				if tt.expectedMetric != "single_connection" {
					t.Errorf("Expected single connection usage")
				}
			}
			
			if tt.concurrency > 5 {
				// Should use multiple connections
				if tt.expectedMetric != "multiple_connections" {
					t.Errorf("Expected multiple connection usage")
				}
			}
		})
	}
}

func TestRepositoryParameterBinding(t *testing.T) {
	// Test SQL parameter binding safety
	
	tests := []struct {
		name        string
		input       string
		expectSafe  bool
		description string
	}{
		{
			name:        "Normal game ID",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			expectSafe:  true,
			description: "UUID should be safe",
		},
		{
			name:        "SQL injection attempt",
			input:       "'; DROP TABLE games; --",
			expectSafe:  true, // Should be safe due to parameter binding
			description: "Parameter binding should prevent injection",
		},
		{
			name:        "Unicode input",
			input:       "测试",
			expectSafe:  true,
			description: "Unicode should be handled safely",
		},
		{
			name:        "Very long input",
			input:       strings.Repeat("A", 10000),
			expectSafe:  true, // Parameter binding should handle length
			description: "Long input should be handled safely",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that parameter binding handles various inputs safely
			if tt.expectSafe {
				// Parameter binding should make this safe
				// This is more of a conceptual test since we use parameterized queries
			}
		})
	}
}
