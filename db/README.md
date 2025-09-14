# Wordle Database

PostgreSQL database setup for the Wordle game application.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Check if it's running
docker-compose ps

# View logs
docker-compose logs postgres
```

### Using Docker directly

```bash
# Build the PostgreSQL image
docker build -f Dockerfile.postgres -t wordle-postgres .

# Run the container
docker run -d \
  --name wordle-postgres \
  -p 5432:5432 \
  -e POSTGRES_DB=wordle \
  -e POSTGRES_USER=wordle_user \
  -e POSTGRES_PASSWORD=wordle_password \
  wordle-postgres
```

## Database Connection

- **Host**: localhost
- **Port**: 5432
- **Database**: wordle
- **Username**: wordle_user
- **Password**: wordle_password

### Connection String
```
postgres://wordle_user:wordle_password@localhost:5432/wordle
```

## Database Schema

### Tables

#### `games`
Stores individual game sessions
- `id` (UUID) - Primary key
- `target_word` (VARCHAR) - The word to guess
- `created_at` (TIMESTAMP) - When the game started
- `completed_at` (TIMESTAMP) - When the game ended
- `is_completed` (BOOLEAN) - Whether the game is finished
- `is_won` (BOOLEAN) - Whether the player won
- `guess_count` (INTEGER) - Number of guesses made
- `max_guesses` (INTEGER) - Maximum allowed guesses (default: 6)

#### `guesses`
Stores individual guesses for each game
- `id` (UUID) - Primary key
- `game_id` (UUID) - Foreign key to games table
- `guess_word` (VARCHAR) - The guessed word
- `guess_number` (INTEGER) - Order of the guess (1-6)
- `result` (JSONB) - Result for each letter (correct/present/absent)
- `created_at` (TIMESTAMP) - When the guess was made

#### `players` (Optional)
Stores player information and statistics
- `id` (UUID) - Primary key
- `username` (VARCHAR) - Player username
- `email` (VARCHAR) - Player email
- `games_played` (INTEGER) - Total games played
- `games_won` (INTEGER) - Total games won
- `current_streak` (INTEGER) - Current winning streak
- `max_streak` (INTEGER) - Maximum winning streak achieved

#### `game_stats` (Optional)
Stores additional game analytics
- `id` (UUID) - Primary key
- `game_id` (UUID) - Foreign key to games table
- `player_id` (UUID) - Foreign key to players table
- `word_difficulty` (FLOAT) - Calculated word difficulty
- `solve_time_seconds` (INTEGER) - Time taken to solve

## Sample Data

The database comes pre-loaded with sample data:
- 3 sample games (completed win, completed loss, in-progress)
- Corresponding guesses with realistic game progression
- 3 sample players with statistics
- Game statistics for completed games

## Connecting from Go

Example Go code to connect to the database:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
    connStr := "postgres://wordle_user:wordle_password@localhost:5432/wordle?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## Management Commands

```bash
# Stop the database
docker-compose down

# Stop and remove all data
docker-compose down -v

# Restart the database
docker-compose restart postgres

# Access PostgreSQL shell
docker-compose exec postgres psql -U wordle_user -d wordle

# Backup database
docker-compose exec postgres pg_dump -U wordle_user wordle > backup.sql

# Restore database
docker-compose exec -T postgres psql -U wordle_user wordle < backup.sql
```

## Development

### Resetting the Database

To reset the database with fresh data:

```bash
docker-compose down -v
docker-compose up -d postgres
```

### Adding New Migrations

Add new SQL files to `db/init/` with sequential numbering:
- `03-add-new-feature.sql`
- `04-update-schema.sql`

Files are executed in alphabetical order during container initialization.
