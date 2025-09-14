package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// GameRepository handles database operations for games
type GameRepository struct {
	db *DB
}

// GuessRepository handles database operations for guesses
type GuessRepository struct {
	db *DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *DB) *GameRepository {
	return &GameRepository{db: db}
}

// NewGuessRepository creates a new guess repository
func NewGuessRepository(db *DB) *GuessRepository {
	return &GuessRepository{db: db}
}

// Game Repository Methods

// CreateGame creates a new game in the database
func (r *GameRepository) CreateGame(targetWord string, maxGuesses int) (*Game, error) {
	query := `
		INSERT INTO games (target_word, max_guesses, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, target_word, created_at, completed_at, is_completed, is_won, guess_count, max_guesses`

	game := &Game{}
	err := r.db.QueryRow(query, targetWord, maxGuesses).Scan(
		&game.ID,
		&game.TargetWord,
		&game.CreatedAt,
		&game.CompletedAt,
		&game.IsCompleted,
		&game.IsWon,
		&game.GuessCount,
		&game.MaxGuesses,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

// GetGame retrieves a game by ID
func (r *GameRepository) GetGame(gameID string) (*Game, error) {
	query := `
		SELECT id, target_word, created_at, completed_at, is_completed, is_won, guess_count, max_guesses
		FROM games
		WHERE id = $1`

	game := &Game{}
	err := r.db.QueryRow(query, gameID).Scan(
		&game.ID,
		&game.TargetWord,
		&game.CreatedAt,
		&game.CompletedAt,
		&game.IsCompleted,
		&game.IsWon,
		&game.GuessCount,
		&game.MaxGuesses,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found: %s", gameID)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

// UpdateGame updates a game in the database
func (r *GameRepository) UpdateGame(game *Game) error {
	query := `
		UPDATE games 
		SET completed_at = $2, is_completed = $3, is_won = $4, guess_count = $5
		WHERE id = $1`

	result, err := r.db.Exec(query,
		game.ID,
		game.CompletedAt,
		game.IsCompleted,
		game.IsWon,
		game.GuessCount,
	)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game not found: %s", game.ID)
	}

	return nil
}

// DeleteGame deletes a game and all associated guesses
func (r *GameRepository) DeleteGame(gameID string) error {
	query := `DELETE FROM games WHERE id = $1`

	result, err := r.db.Exec(query, gameID)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game not found: %s", gameID)
	}

	return nil
}

// GetGameWithGuesses retrieves a game with all its guesses
func (r *GameRepository) GetGameWithGuesses(gameID string) (*GameWithGuesses, error) {
	game, err := r.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	guessRepo := NewGuessRepository(r.db)
	guesses, err := guessRepo.GetGuessesByGameID(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guesses: %w", err)
	}

	return &GameWithGuesses{
		Game:    *game,
		Guesses: guesses,
	}, nil
}

// GetRecentGames gets the most recent games
func (r *GameRepository) GetRecentGames(limit int) ([]Game, error) {
	query := `
		SELECT id, target_word, created_at, completed_at, is_completed, is_won, guess_count, max_guesses
		FROM games
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent games: %w", err)
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		err := rows.Scan(
			&game.ID,
			&game.TargetWord,
			&game.CreatedAt,
			&game.CompletedAt,
			&game.IsCompleted,
			&game.IsWon,
			&game.GuessCount,
			&game.MaxGuesses,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// Guess Repository Methods

// CreateGuess creates a new guess in the database
func (r *GuessRepository) CreateGuess(gameID, guessWord string, guessNumber int, result GuessResult) (*Guess, error) {
	query := `
		INSERT INTO guesses (game_id, guess_word, guess_number, result, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, game_id, guess_word, guess_number, result, created_at`

	guess := &Guess{}
	err := r.db.QueryRow(query, gameID, guessWord, guessNumber, result).Scan(
		&guess.ID,
		&guess.GameID,
		&guess.GuessWord,
		&guess.GuessNumber,
		&guess.Result,
		&guess.CreatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return nil, fmt.Errorf("guess number %d already exists for game %s", guessNumber, gameID)
			}
		}
		return nil, fmt.Errorf("failed to create guess: %w", err)
	}

	return guess, nil
}

// GetGuess retrieves a guess by ID
func (r *GuessRepository) GetGuess(guessID string) (*Guess, error) {
	query := `
		SELECT id, game_id, guess_word, guess_number, result, created_at
		FROM guesses
		WHERE id = $1`

	guess := &Guess{}
	err := r.db.QueryRow(query, guessID).Scan(
		&guess.ID,
		&guess.GameID,
		&guess.GuessWord,
		&guess.GuessNumber,
		&guess.Result,
		&guess.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("guess not found: %s", guessID)
		}
		return nil, fmt.Errorf("failed to get guess: %w", err)
	}

	return guess, nil
}

// GetGuessesByGameID retrieves all guesses for a game, ordered by guess number
func (r *GuessRepository) GetGuessesByGameID(gameID string) ([]Guess, error) {
	query := `
		SELECT id, game_id, guess_word, guess_number, result, created_at
		FROM guesses
		WHERE game_id = $1
		ORDER BY guess_number ASC`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guesses: %w", err)
	}
	defer rows.Close()

	var guesses []Guess
	for rows.Next() {
		var guess Guess
		err := rows.Scan(
			&guess.ID,
			&guess.GameID,
			&guess.GuessWord,
			&guess.GuessNumber,
			&guess.Result,
			&guess.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guess: %w", err)
		}
		guesses = append(guesses, guess)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating guesses: %w", err)
	}

	return guesses, nil
}

// DeleteGuess deletes a guess
func (r *GuessRepository) DeleteGuess(guessID string) error {
	query := `DELETE FROM guesses WHERE id = $1`

	result, err := r.db.Exec(query, guessID)
	if err != nil {
		return fmt.Errorf("failed to delete guess: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("guess not found: %s", guessID)
	}

	return nil
}

// GetLatestGuess gets the most recent guess for a game
func (r *GuessRepository) GetLatestGuess(gameID string) (*Guess, error) {
	query := `
		SELECT id, game_id, guess_word, guess_number, result, created_at
		FROM guesses
		WHERE game_id = $1
		ORDER BY guess_number DESC
		LIMIT 1`

	guess := &Guess{}
	err := r.db.QueryRow(query, gameID).Scan(
		&guess.ID,
		&guess.GameID,
		&guess.GuessWord,
		&guess.GuessNumber,
		&guess.Result,
		&guess.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no guesses found for game: %s", gameID)
		}
		return nil, fmt.Errorf("failed to get latest guess: %w", err)
	}

	return guess, nil
}
