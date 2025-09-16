package main

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalEnvVars := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"PORT":        os.Getenv("PORT"),
		"MAX_GUESSES": os.Getenv("MAX_GUESSES"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalEnvVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test with default values
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("MAX_GUESSES")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not return error: %v", err)
	}

	// Test default values
	if config.Database.Host != "localhost" {
		t.Errorf("Expected default DB host 'localhost', got '%s'", config.Database.Host)
	}
	if config.Database.Port != 5432 {
		t.Errorf("Expected default DB port 5432, got %d", config.Database.Port)
	}
	if config.Database.Name != "wordle" {
		t.Errorf("Expected default DB name 'wordle', got '%s'", config.Database.Name)
	}
	if config.Server.Port != 8080 {
		t.Errorf("Expected default server port 8080, got %d", config.Server.Port)
	}
	if config.Game.MaxGuesses != 6 {
		t.Errorf("Expected default max guesses 6, got %d", config.Game.MaxGuesses)
	}
	if config.Game.WordLength != 5 {
		t.Errorf("Expected default word length 5, got %d", config.Game.WordLength)
	}
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Set custom environment variables
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "custom_db")
	os.Setenv("DB_USER", "custom_user")
	os.Setenv("DB_PASSWORD", "custom_pass")
	os.Setenv("PORT", "9000")
	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("MAX_GUESSES", "10")
	os.Setenv("WORD_LENGTH", "7")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("MAX_GUESSES")
		os.Unsetenv("WORD_LENGTH")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not return error: %v", err)
	}

	// Test custom values
	if config.Database.Host != "custom-host" {
		t.Errorf("Expected DB host 'custom-host', got '%s'", config.Database.Host)
	}
	if config.Database.Port != 3306 {
		t.Errorf("Expected DB port 3306, got %d", config.Database.Port)
	}
	if config.Database.Name != "custom_db" {
		t.Errorf("Expected DB name 'custom_db', got '%s'", config.Database.Name)
	}
	if config.Database.User != "custom_user" {
		t.Errorf("Expected DB user 'custom_user', got '%s'", config.Database.User)
	}
	if config.Database.Password != "custom_pass" {
		t.Errorf("Expected DB password 'custom_pass', got '%s'", config.Database.Password)
	}
	if config.Server.Port != 9000 {
		t.Errorf("Expected server port 9000, got %d", config.Server.Port)
	}
	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected server host '0.0.0.0', got '%s'", config.Server.Host)
	}
	if config.Game.MaxGuesses != 10 {
		t.Errorf("Expected max guesses 10, got %d", config.Game.MaxGuesses)
	}
	if config.Game.WordLength != 7 {
		t.Errorf("Expected word length 7, got %d", config.Game.WordLength)
	}
}

func TestDatabaseConfigConnectionString(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "testhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	expected := "host=testhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	actual := config.ConnectionString()

	if actual != expected {
		t.Errorf("Expected connection string '%s', got '%s'", expected, actual)
	}
}

func TestDatabaseConfigDatabaseURL(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "testhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	expected := "postgres://testuser:testpass@testhost:5432/testdb?sslmode=disable"
	actual := config.DatabaseURL()

	if actual != expected {
		t.Errorf("Expected database URL '%s', got '%s'", expected, actual)
	}
}

func TestServerConfigAddress(t *testing.T) {
	config := &ServerConfig{
		Host: "localhost",
		Port: 8080,
	}

	expected := "localhost:8080"
	actual := config.Address()

	if actual != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, actual)
	}
}

func TestGetEnvString(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_ENV_STRING", "test_value")
	defer os.Unsetenv("TEST_ENV_STRING")

	result := getEnvString("TEST_ENV_STRING", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with non-existing env var
	result = getEnvString("NON_EXISTING_ENV", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid int env var
	os.Setenv("TEST_ENV_INT", "42")
	defer os.Unsetenv("TEST_ENV_INT")

	result := getEnvInt("TEST_ENV_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid int env var
	os.Setenv("TEST_ENV_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_ENV_INVALID_INT")

	result = getEnvInt("TEST_ENV_INVALID_INT", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}

	// Test with non-existing env var
	result = getEnvInt("NON_EXISTING_ENV", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}
}

func TestGetEnvDuration(t *testing.T) {
	// Test with valid duration env var
	os.Setenv("TEST_ENV_DURATION", "30m")
	defer os.Unsetenv("TEST_ENV_DURATION")

	result := getEnvDuration("TEST_ENV_DURATION", "1h")
	expected := 30 * time.Minute
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test with invalid duration env var
	os.Setenv("TEST_ENV_INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("TEST_ENV_INVALID_DURATION")

	result = getEnvDuration("TEST_ENV_INVALID_DURATION", "1h")
	expected = time.Hour
	if result != expected {
		t.Errorf("Expected default value %v, got %v", expected, result)
	}

	// Test with non-existing env var
	result = getEnvDuration("NON_EXISTING_ENV", "1h")
	expected = time.Hour
	if result != expected {
		t.Errorf("Expected default value %v, got %v", expected, result)
	}

	// Test with invalid default duration (should fallback to 1 hour)
	result = getEnvDuration("NON_EXISTING_ENV", "invalid_default")
	expected = time.Hour
	if result != expected {
		t.Errorf("Expected fallback value %v, got %v", expected, result)
	}
}
