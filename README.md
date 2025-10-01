# Wordle Game

A full-stack Wordle implementation with a Go backend API and React frontend client.

## 🎮 Features

### ✅ Complete Game Implementation
- **Interactive React Client** - Modern web interface with real-time gameplay
- **Go REST API Server** - Robust backend with PostgreSQL database
- **Two-Tier Word System** - Smart word selection for better gameplay
- **Used Letters Display** - Visual alphabet grid showing letter status
- **Game State Management** - Persistent games with guess history
- **Statistics Tracking** - Game stats and recent games history

### 🎯 Game Features
- 6 attempts to guess a 5-letter word
- Color-coded feedback (Green/Yellow/Gray)
- Real-time validation of guesses
- Visual alphabet grid showing used letters
- Game statistics and history
- Responsive design for mobile and desktop

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+** for the backend server
- **Node.js 14+** for the React client
- **PostgreSQL** for the database

### 1. Start the Database
```bash
# Using Docker (recommended)
docker run --name wordle-postgres -e POSTGRES_DB=wordle -e POSTGRES_USER=wordle_user -e POSTGRES_PASSWORD=wordle_password -p 5432:5432 -d postgres:13

# Or use your local PostgreSQL installation
createdb wordle
psql -d wordle -c "CREATE USER wordle_user WITH PASSWORD 'wordle_password';"
psql -d wordle -c "GRANT ALL PRIVILEGES ON DATABASE wordle TO wordle_user;"
```

### 2. Start the Go Server
```bash
cd server
go run .
```
Server starts on `http://localhost:8080`

### 3. Start the React Client
```bash
cd client
npm install
npm start
```
Client starts on `http://localhost:3000`

### 4. Play Wordle!
Open `http://localhost:3000` in your browser and start playing!

## 🏗️ Architecture

### Backend (Go)
- **REST API** with endpoints for game management
- **PostgreSQL Database** for persistent storage
- **Two-Tier Word System**:
  - 14,855 validation words for guess checking
  - 527 common target words for better gameplay
- **Comprehensive Testing** with unit and integration tests

### Frontend (React)
- **Modern React 17** with hooks and functional components
- **Real-time Game Board** with color-coded tiles
- **Used Letters Display** showing alphabet status
- **Game Statistics** and recent games history
- **Responsive Design** for all screen sizes

### Database Schema
```sql
-- Games table
CREATE TABLE games (
    id UUID PRIMARY KEY,
    target_word VARCHAR(5) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    is_completed BOOLEAN DEFAULT FALSE,
    is_won BOOLEAN DEFAULT FALSE,
    guess_count INTEGER DEFAULT 0,
    max_guesses INTEGER DEFAULT 6
);

-- Guesses table
CREATE TABLE guesses (
    id UUID PRIMARY KEY,
    game_id UUID REFERENCES games(id),
    guess_word VARCHAR(5) NOT NULL,
    guess_number INTEGER NOT NULL,
    result JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/games` | Create a new game |
| `GET` | `/api/games/{id}` | Get game state with guesses |
| `POST` | `/api/games/{id}` | Make a guess |
| `DELETE` | `/api/games/{id}` | Delete a game |
| `GET` | `/api/games` | Get recent games |
| `GET` | `/api/stats` | Get game statistics |
| `GET` | `/health` | Health check |

### Example API Usage

```bash
# Create a new game
curl -X POST http://localhost:8080/api/games

# Make a guess
curl -X POST http://localhost:8080/api/games/{game-id} \
  -H "Content-Type: application/json" \
  -d '{"guess_word": "CRANE"}'

# Get game state
curl http://localhost:8080/api/games/{game-id}
```

## 🎨 Word System

### Validation Words (14,855 words)
- **File**: `server/valid-wordle-words.txt`
- **Purpose**: Validate all player guesses
- **Examples**: `aahed`, `zymes`, `crwth`, etc.

### Target Words (527 words)
- **File**: `server/common-target-words.txt`
- **Purpose**: Generate game targets (more enjoyable)
- **Examples**: `about`, `world`, `house`, `money`, etc.

This two-tier system ensures players get familiar, common words as targets while still allowing any valid English word as a guess.

## 🧪 Development

### Running Tests
```bash
# Backend tests
cd server
go test -v

# Test specific functionality
go test -v -run TestWordList
go test -v -run TestGameService
```

### Building for Production
```bash
# Build Go server
cd server
go build -o wordle-server .

# Build React client
cd client
npm run build
```

### Database Management
```bash
# Run database migrations
cd server
go run . --migrate

# Reset database (development only)
psql -d wordle -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
```

## 🎯 Game Rules

1. **Objective**: Guess the 5-letter word in 6 attempts
2. **Feedback Colors**:
   - 🟩 **Green**: Correct letter in correct position
   - 🟨 **Yellow**: Correct letter in wrong position  
   - ⬜ **Gray**: Letter not in the word
3. **Used Letters**: Alphabet grid shows status of all guessed letters
4. **Valid Words**: Only real English words are accepted

## 📊 Features in Detail

### Used Letters Display
- Visual alphabet grid (A-Z)
- Color-coded based on guess results
- Smart status priority (correct > present > absent)
- Responsive layout for mobile devices

### Game Statistics
- Total games played
- Win percentage
- Average number of guesses
- Recent games history

### Responsive Design
- Desktop: Full-size game board and alphabet
- Mobile: Optimized layouts and touch-friendly controls
- Dark theme matching original Wordle aesthetic

## 🔧 Configuration

### Environment Variables
```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=wordle
DB_USER=wordle_user
DB_PASSWORD=wordle_password

# Server configuration
SERVER_HOST=localhost
SERVER_PORT=8080
```

### Client Configuration
The React client automatically proxies API requests to `http://localhost:8080`. Update `client/package.json` if your backend runs on a different port.

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Run tests**: `go test ./server/...` and `npm test` in client
4. **Commit changes**: `git commit -m 'Add amazing feature'`
5. **Push to branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

### Code Standards
- **Go**: Follow standard Go conventions and run `go fmt`
- **React**: Use ESLint configuration and Prettier formatting
- **Tests**: Add tests for new functionality
- **Documentation**: Update README for new features

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎉 Acknowledgments

- Original Wordle game by Josh Wardle
- Word list sourced from various English dictionaries
- Built with Go, React, and PostgreSQL