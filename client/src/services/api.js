// API service for interacting with the Wordle backend
const API_BASE_URL = process.env.REACT_APP_API_URL || '';

class WordleAPI {
  async makeRequest(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || `HTTP error! status: ${response.status}`);
      }

      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Create a new game
  async createGame(maxGuesses = 6) {
    return this.makeRequest('/api/games', {
      method: 'POST',
      body: JSON.stringify({ max_guesses: maxGuesses }),
    });
  }

  // Get game state with guesses
  async getGame(gameId) {
    return this.makeRequest(`/api/games/${gameId}`);
  }

  // Make a guess
  async makeGuess(gameId, guessWord) {
    return this.makeRequest(`/api/games/${gameId}`, {
      method: 'POST',
      body: JSON.stringify({ guess_word: guessWord }),
    });
  }

  // Get recent games
  async getRecentGames() {
    return this.makeRequest('/api/games');
  }

  // Get game statistics
  async getStats() {
    return this.makeRequest('/api/stats');
  }

  // Delete a game
  async deleteGame(gameId) {
    return this.makeRequest(`/api/games/${gameId}`, {
      method: 'DELETE',
    });
  }

  // Health check
  async healthCheck() {
    return this.makeRequest('/health');
  }
}

const wordleAPI = new WordleAPI();
export default wordleAPI;
