import React, { useState, useEffect } from 'react';
import api from '../services/api';

const GameStats = () => {
  const [stats, setStats] = useState(null);
  const [recentGames, setRecentGames] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadStats();
    loadRecentGames();
  }, []);

  const loadStats = async () => {
    setLoading(true);
    try {
      const statsData = await api.getStats();
      setStats(statsData);
    } catch (err) {
      setError(`Failed to load stats: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const loadRecentGames = async () => {
    try {
      const gamesData = await api.getRecentGames();
      setRecentGames(gamesData.games || []);
    } catch (err) {
      console.error('Failed to load recent games:', err);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  const getGameStatus = (game) => {
    if (game.is_won) return 'won';
    if (game.is_completed) return 'lost';
    return 'in-progress';
  };

  const getGameStatusText = (game) => {
    if (game.is_won) return `Won in ${game.guess_count} guesses`;
    if (game.is_completed) return 'Lost';
    return `${game.guess_count}/${game.max_guesses} guesses`;
  };

  return (
    <div>
      {loading && <div className="loading">Loading stats...</div>}
      {error && <div className="error-message">{error}</div>}

      {stats && (
        <div className="stats">
          <h3>Game Statistics</h3>
          <div>Total Games: {stats.total_games || 0}</div>
          <div>Games Won: {stats.games_won || 0}</div>
          <div>Win Rate: {stats.win_rate ? `${stats.win_rate.toFixed(1)}%` : '0%'}</div>
          <div>Average Guesses: {stats.average_guesses ? stats.average_guesses.toFixed(1) : 'N/A'}</div>
        </div>
      )}

      {recentGames.length > 0 && (
        <div className="recent-games">
          <h3>Recent Games</h3>
          {recentGames.map((game) => (
            <div key={game.id} className="game-item">
              <div className="game-item-info">
                <div>Game {game.id.slice(0, 8)}...</div>
                <div>{formatDate(game.created_at)}</div>
              </div>
              <div className={`game-item-status ${getGameStatus(game)}`}>
                {getGameStatusText(game)}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default GameStats;
