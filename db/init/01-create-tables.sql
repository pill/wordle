-- Create tables for Wordle game

-- Games table to store individual game sessions
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_word VARCHAR(5) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    is_completed BOOLEAN DEFAULT FALSE,
    is_won BOOLEAN DEFAULT FALSE,
    guess_count INTEGER DEFAULT 0,
    max_guesses INTEGER DEFAULT 6
);

-- Guesses table to store individual guesses for each game
CREATE TABLE IF NOT EXISTS guesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    guess_word VARCHAR(5) NOT NULL,
    guess_number INTEGER NOT NULL,
    result JSONB NOT NULL, -- Store the result as JSON (correct, present, absent for each letter)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(game_id, guess_number)
);

-- Players table (optional, for future user management)
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    games_played INTEGER DEFAULT 0,
    games_won INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    max_streak INTEGER DEFAULT 0
);

-- Game statistics (optional, for analytics)
CREATE TABLE IF NOT EXISTS game_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id UUID REFERENCES players(id) ON DELETE SET NULL,
    word_difficulty FLOAT, -- Could be calculated based on word frequency
    solve_time_seconds INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_games_created_at ON games(created_at);
CREATE INDEX IF NOT EXISTS idx_games_target_word ON games(target_word);
CREATE INDEX IF NOT EXISTS idx_guesses_game_id ON guesses(game_id);
CREATE INDEX IF NOT EXISTS idx_guesses_created_at ON guesses(created_at);
CREATE INDEX IF NOT EXISTS idx_players_username ON players(username);
CREATE INDEX IF NOT EXISTS idx_game_stats_game_id ON game_stats(game_id);
CREATE INDEX IF NOT EXISTS idx_game_stats_player_id ON game_stats(player_id);
