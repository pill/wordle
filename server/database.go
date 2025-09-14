package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the database connection with additional functionality
type DB struct {
	*sql.DB
	config *DatabaseConfig
}

// NewDB creates a new database connection with proper configuration
func NewDB(config *DatabaseConfig) (*DB, error) {
	// Open database connection
	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database: %s:%d/%s", config.Host, config.Port, config.Name)

	return &DB{
		DB:     db,
		config: config,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	log.Println("Closing database connection")
	return db.DB.Close()
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// Stats returns connection pool statistics
func (db *DB) Stats() sql.DBStats {
	return db.DB.Stats()
}

// Config returns the database configuration
func (db *DB) Config() *DatabaseConfig {
	return db.config
}

// HealthCheck performs a comprehensive health check of the database
func (db *DB) HealthCheck() error {
	// Check basic connectivity
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("query test failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	return nil
}

// BeginTx starts a new transaction with the given options
func (db *DB) BeginTx(opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(nil, opts)
}

// ExecContext executes a query without returning any rows with logging
func (db *DB) ExecWithLog(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.Exec(query, args...)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Query failed (took %v): %s, args: %v, error: %v", duration, query, args, err)
	} else {
		log.Printf("Query executed (took %v): %s", duration, query)
	}

	return result, err
}

// QueryWithLog executes a query that returns rows with logging
func (db *DB) QueryWithLog(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.Query(query, args...)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Query failed (took %v): %s, args: %v, error: %v", duration, query, args, err)
	} else {
		log.Printf("Query executed (took %v): %s", duration, query)
	}

	return rows, err
}

// QueryRowWithLog executes a query that returns at most one row with logging
func (db *DB) QueryRowWithLog(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := db.DB.QueryRow(query, args...)
	duration := time.Since(start)

	log.Printf("Query executed (took %v): %s", duration, query)
	return row
}

// Migrate runs database migrations (placeholder for future migration system)
func (db *DB) Migrate() error {
	// This is a placeholder for a more sophisticated migration system
	// For now, we'll just verify that the required tables exist
	
	tables := []string{"games", "guesses", "players", "game_stats"}
	
	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)`
		
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}
		
		if !exists {
			return fmt.Errorf("required table %s does not exist", table)
		}
		
		log.Printf("Table %s exists", table)
	}
	
	log.Println("All required tables exist")
	return nil
}

// LogConnectionStats logs current connection pool statistics
func (db *DB) LogConnectionStats() {
	stats := db.Stats()
	log.Printf("DB Connection Stats - Open: %d, InUse: %d, Idle: %d, Wait: %d, MaxOpen: %d",
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.MaxOpenConnections,
	)
}
