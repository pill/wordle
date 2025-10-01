package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Global variables
var (
	gameService *GameService
	config      *Config
)

func main() {
	// Load configuration
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize word list
	wordList, err := NewWordList("")
	if err != nil {
		log.Fatalf("Failed to initialize word list: %v", err)
	}

	// Initialize database connection
	db, err := NewDB(&config.Database)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Running in demo mode without database...")
		runWordListDemo(wordList)
		return
	}
	defer db.Close()

	// Run database migrations/checks
	if err := db.Migrate(); err != nil {
		log.Printf("Warning: Database migration check failed: %v", err)
		log.Println("Running in demo mode without database...")
		runWordListDemo(wordList)
		return
	}

	// Initialize game service
	gameService = NewGameService(db, wordList, &config.Game)

	// Setup HTTP handlers
	setupRoutes()

	// Start server
	address := config.Server.Address()
	log.Printf("Wordle API server starting on %s...", address)
	log.Printf("Database connected: %s", config.Database.DatabaseURL())
	log.Printf("Word lists loaded: %d validation words, %d target words", wordList.Size(), wordList.TargetWordsSize())
	
	log.Fatal(http.ListenAndServe(address, nil))
}

func setupRoutes() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/games", gamesHandler)
	http.HandleFunc("/api/games/", gameHandler) // for /api/games/{id}
	http.HandleFunc("/api/stats", statsHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to the Wordle API!",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"POST /api/games":      "Create a new game",
			"GET /api/games/{id}":  "Get game state",
			"POST /api/games/{id}": "Make a guess",
			"GET /api/stats":       "Get game statistics",
			"GET /health":          "Health check",
		},
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"database":  "connected",
	}
	writeJSONResponse(w, http.StatusOK, status)
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createGameHandler(w, r)
	case http.MethodGet:
		getRecentGamesHandler(w, r)
	default:
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/games/")
	gameID := strings.Split(path, "/")[0]
	
	if gameID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Game ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		getGameHandler(w, r, gameID)
	case http.MethodPost:
		makeGuessHandler(w, r, gameID)
	case http.MethodDelete:
		deleteGameHandler(w, r, gameID)
	default:
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	game, err := gameService.CreateNewGame()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create game: %v", err))
		return
	}

	response := GameResponse{
		Game:    *game,
		Message: fmt.Sprintf("New game created! You have %d guesses to find the word.", game.MaxGuesses),
	}

	writeJSONResponse(w, http.StatusCreated, response)
}

func getGameHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	gameWithGuesses, err := gameService.GetGameWithGuesses(gameID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get game: %v", err))
		}
		return
	}

	response := GameResponse{
		Game:    gameWithGuesses.Game,
		Guesses: gameWithGuesses.Guesses,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func makeGuessHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	var request MakeGuessRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if request.GuessWord == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Guess word is required")
		return
	}

	response, err := gameService.MakeGuess(gameID, request.GuessWord)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Game not found")
		} else if strings.Contains(err.Error(), "not a valid word") || 
		          strings.Contains(err.Error(), "must be") ||
		          strings.Contains(err.Error(), "already completed") ||
		          strings.Contains(err.Error(), "no remaining") {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process guess: %v", err))
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func deleteGameHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	err := gameService.DeleteGame(gameID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Game not found")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete game: %v", err))
		}
		return
	}

	response := map[string]string{
		"message": "Game deleted successfully",
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func getRecentGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := gameService.GetRecentGames(10)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get recent games: %v", err))
		return
	}

	response := map[string]interface{}{
		"games": games,
		"count": len(games),
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := gameService.GetGameStats()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get stats: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, stats)
}

// Helper functions

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}
	writeJSONResponse(w, statusCode, response)
}

// runWordListDemo runs the original word list demo when database is not available
func runWordListDemo(wordList *WordList) {
	fmt.Println("=== WordList Demo Mode ===")
	fmt.Printf("Total words loaded: %d\n", wordList.Size())

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Random word: %s\n", wordList.RandomWord())
	fmt.Printf("Random word: %s\n", wordList.RandomWord())
	fmt.Printf("Random word: %s\n", wordList.RandomWord())

	fmt.Println("\n=== Five Letter Words ===")
	fiveLetterWords := wordList.FiveLetterWords()
	fmt.Printf("Number of five-letter words: %d\n", len(fiveLetterWords))

	// Get sample five-letter words
	sampleWords := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(fiveLetterWords); i++ {
		idx := rand.Intn(len(fiveLetterWords))
		sampleWords = append(sampleWords, fiveLetterWords[idx])
	}
	fmt.Printf("Sample five-letter words: %s\n", strings.Join(sampleWords, ", "))

	fmt.Println("\n=== Word Validation ===")
	testWords := []string{"hello", "world", "apple", "xyzzy", "abask", "aahed"}
	for _, word := range testWords {
		valid := wordList.Contains(word)
		status := "not valid"
		if valid {
			status = "valid"
		}
		fmt.Printf("'%s' is %s\n", word, status)
	}

	fmt.Println("\n=== Word Length Distribution ===")
	for length := 3; length <= 8; length++ {
		count := len(wordList.WordsOfLength(length))
		fmt.Printf("%d-letter words: %d\n", length, count)
	}

	fmt.Println("\n=== Sample Words by Length ===")
	for length := 3; length <= 8; length++ {
		words := wordList.WordsOfLength(length)
		if len(words) > 0 {
			sampleSize := 5
			if len(words) < sampleSize {
				sampleSize = len(words)
			}

			sample := make([]string, 0, sampleSize)
			for i := 0; i < sampleSize; i++ {
				idx := rand.Intn(len(words))
				sample = append(sample, words[idx])
			}
			fmt.Printf("%d-letter words sample: %s\n", length, strings.Join(sample, ", "))
		}
	}
}
