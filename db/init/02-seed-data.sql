-- Seed data for Wordle database

-- Insert some sample games for testing
INSERT INTO games (id, target_word, created_at, completed_at, is_completed, is_won, guess_count) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'HELLO', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', true, true, 4),
    ('550e8400-e29b-41d4-a716-446655440002', 'WORLD', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', true, false, 6),
    ('550e8400-e29b-41d4-a716-446655440003', 'APPLE', NOW() - INTERVAL '1 hour', NULL, false, false, 2);

-- Insert corresponding guesses
INSERT INTO guesses (game_id, guess_word, guess_number, result) VALUES
    -- Game 1 (HELLO) - Won in 4 guesses
    ('550e8400-e29b-41d4-a716-446655440001', 'CRANE', 1, '[{"letter":"C","status":"absent"},{"letter":"R","status":"absent"},{"letter":"A","status":"absent"},{"letter":"N","status":"absent"},{"letter":"E","status":"present"}]'),
    ('550e8400-e29b-41d4-a716-446655440001', 'SLIME', 2, '[{"letter":"S","status":"absent"},{"letter":"L","status":"correct"},{"letter":"I","status":"absent"},{"letter":"M","status":"absent"},{"letter":"E","status":"present"}]'),
    ('550e8400-e29b-41d4-a716-446655440001', 'HELLO', 3, '[{"letter":"H","status":"correct"},{"letter":"E","status":"correct"},{"letter":"L","status":"correct"},{"letter":"L","status":"correct"},{"letter":"O","status":"correct"}]'),
    
    -- Game 2 (WORLD) - Lost after 6 guesses
    ('550e8400-e29b-41d4-a716-446655440002', 'CRANE', 1, '[{"letter":"C","status":"absent"},{"letter":"R","status":"correct"},{"letter":"A","status":"absent"},{"letter":"N","status":"absent"},{"letter":"E","status":"absent"}]'),
    ('550e8400-e29b-41d4-a716-446655440002', 'GROUT', 2, '[{"letter":"G","status":"absent"},{"letter":"R","status":"correct"},{"letter":"O","status":"present"},{"letter":"U","status":"absent"},{"letter":"T","status":"absent"}]'),
    ('550e8400-e29b-41d4-a716-446655440002', 'PRODS', 3, '[{"letter":"P","status":"absent"},{"letter":"R","status":"correct"},{"letter":"O","status":"present"},{"letter":"D","status":"present"},{"letter":"S","status":"absent"}]'),
    ('550e8400-e29b-41d4-a716-446655440002', 'LORDS', 4, '[{"letter":"L","status":"present"},{"letter":"O","status":"present"},{"letter":"R","status":"correct"},{"letter":"D","status":"present"},{"letter":"S","status":"absent"}]'),
    ('550e8400-e29b-41d4-a716-446655440002', 'DROLL', 5, '[{"letter":"D","status":"present"},{"letter":"R","status":"correct"},{"letter":"O","status":"present"},{"letter":"L","status":"present"},{"letter":"L","status":"present"}]'),
    ('550e8400-e29b-41d4-a716-446655440002', 'WRONG', 6, '[{"letter":"W","status":"correct"},{"letter":"R","status":"correct"},{"letter":"O","status":"correct"},{"letter":"N","status":"absent"},{"letter":"G","status":"absent"}]'),
    
    -- Game 3 (APPLE) - In progress
    ('550e8400-e29b-41d4-a716-446655440003', 'CRANE', 1, '[{"letter":"C","status":"absent"},{"letter":"R","status":"absent"},{"letter":"A","status":"present"},{"letter":"N","status":"absent"},{"letter":"E","status":"present"}]'),
    ('550e8400-e29b-41d4-a716-446655440003', 'PLEAS', 2, '[{"letter":"P","status":"correct"},{"letter":"L","status":"correct"},{"letter":"E","status":"present"},{"letter":"A","status":"present"},{"letter":"S","status":"absent"}]');

-- Insert sample players
INSERT INTO players (id, username, email, games_played, games_won, current_streak, max_streak) VALUES
    ('660e8400-e29b-41d4-a716-446655440001', 'wordmaster', 'wordmaster@example.com', 25, 20, 5, 12),
    ('660e8400-e29b-41d4-a716-446655440002', 'puzzler', 'puzzler@example.com', 15, 10, 0, 8),
    ('660e8400-e29b-41d4-a716-446655440003', 'newbie', 'newbie@example.com', 3, 1, 1, 1);

-- Insert game statistics
INSERT INTO game_stats (game_id, player_id, word_difficulty, solve_time_seconds) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', 0.6, 180),
    ('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440002', 0.8, 420);
