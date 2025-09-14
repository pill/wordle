# Database Setup Guide

This guide will help you set up and use the PostgreSQL database for the Wordle application.

## Quick Start

### 1. Start the Database

```bash
# Start PostgreSQL with Docker Compose
docker-compose up -d postgres

# Verify it's running
docker-compose ps

# Check logs
docker-compose logs postgres
```

### 2. Run the Application

```bash
# From the server directory
cd server
go run .
```

The application will automatically:
- Connect to the database
- Verify required tables exist
- Start the API server on port 8080

## API Endpoints

### Game Management

#### Create a New Game
```bash
curl -X POST http://localhost:8080/api/games
```

Response:
```json
{
  "game": {
    "id": "uuid-here",
    "target_word": "HIDDEN",
    "created_at": "2025-09-14T10:47:31Z",
    "is_completed": false,
    "is_won": false,
    "guess_count": 0,
    "max_guesses": 6
  },
  "message": "New game created! You have 6 guesses to find the word."
}
```

#### Get Game State
```bash
curl http://localhost:8080/api/games/{game-id}
```

#### Make a Guess
```bash
curl -X POST http://localhost:8080/api/games/{game-id} \
  -H "Content-Type: application/json" \
  -d '{"guess_word": "HELLO"}'
```

Response:
```json
{
  "game": {
    "id": "uuid-here",
    "target_word": "HIDDEN",
    "created_at": "2025-09-14T10:47:31Z",
    "is_completed": false,
    "is_won": false,
    "guess_count": 1,
    "max_guesses": 6
  },
  "guesses": [
    {
      "id": "guess-uuid",
      "game_id": "game-uuid",
      "guess_word": "HELLO",
      "guess_number": 1,
      "result": [
        {"letter": "H", "status": "absent"},
        {"letter": "E", "status": "present"},
        {"letter": "L", "status": "correct"},
        {"letter": "L", "status": "absent"},
        {"letter": "O", "status": "correct"}
      ],
      "created_at": "2025-09-14T10:47:35Z"
    }
  ],
  "message": "Good guess! 5 guess(es) remaining"
}
```

### Other Endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Game Statistics
```bash
curl http://localhost:8080/api/stats
```

#### Recent Games
```bash
curl http://localhost:8080/api/games
```

## Environment Configuration

Create a `.env` file in the project root (copy from `env.example`):

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=wordle
DB_USER=wordle_user
DB_PASSWORD=wordle_password

# Server Configuration
PORT=8080
HOST=localhost

# Game Configuration
MAX_GUESSES=6
WORD_LENGTH=5
```

## Database Management Scripts

Use the included database management script:

```bash
# Start database
./scripts/db.sh start

# Stop database
./scripts/db.sh stop

# Reset database (deletes all data)
./scripts/db.sh reset

# Connect to PostgreSQL shell
./scripts/db.sh shell

# Create backup
./scripts/db.sh backup

# Restore from backup
./scripts/db.sh restore backup_file.sql

# View logs
./scripts/db.sh logs

# Check status
./scripts/db.sh status
```

## Game Logic

### Guess Evaluation

When you make a guess, each letter is evaluated against the target word:

- **Correct** (🟩): Letter is in the correct position
- **Present** (🟨): Letter is in the word but wrong position  
- **Absent** (⬜): Letter is not in the word

### Game Rules

- 6 maximum guesses per game
- 5-letter words only
- Only valid dictionary words accepted
- Case-insensitive input
- Game ends when word is guessed or max guesses reached

## Database Schema

### Tables

- **games**: Game sessions with target words and completion status
- **guesses**: Individual guesses with results for each letter
- **players**: Player information and statistics (optional)
- **game_stats**: Additional game analytics (optional)

### Sample Queries

```sql
-- Get all games
SELECT * FROM games ORDER BY created_at DESC;

-- Get game with guesses
SELECT g.*, gu.guess_word, gu.guess_number, gu.result 
FROM games g 
LEFT JOIN guesses gu ON g.id = gu.game_id 
WHERE g.id = 'game-id'
ORDER BY gu.guess_number;

-- Win rate statistics
SELECT 
  COUNT(*) as total_games,
  COUNT(*) FILTER (WHERE is_won = true) as won_games,
  ROUND(
    COUNT(*) FILTER (WHERE is_won = true) * 100.0 / COUNT(*), 
    2
  ) as win_rate
FROM games 
WHERE is_completed = true;
```

## Troubleshooting

### Database Connection Issues

1. **Check if PostgreSQL is running:**
   ```bash
   docker-compose ps postgres
   ```

2. **Check logs:**
   ```bash
   docker-compose logs postgres
   ```

3. **Reset database:**
   ```bash
   ./scripts/db.sh reset
   ```

### Application Issues

1. **If database tables don't exist:**
   - The app will run in demo mode without database features
   - Check that Docker containers started properly
   - Verify initialization scripts in `db/init/` ran correctly

2. **Port conflicts:**
   - Change `PORT` in `.env` file
   - Or change port mapping in `docker-compose.yml`

3. **Word list not loading:**
   - Ensure `valid-wordle-words.txt` exists in server directory
   - Check file permissions

## Performance Notes

- Database connection pool configured for development (25 max connections)
- Indexes on frequently queried columns (game_id, created_at, etc.)
- JSON storage for guess results enables flexible querying
- Sample data included for immediate testing

## Security Notes

- Default passwords are for development only
- Use strong passwords in production
- Consider SSL/TLS for database connections in production
- API endpoints have basic validation but no authentication yet

