import React, { useState, useEffect } from 'react';
import api from '../services/api';

const WordleGame = () => {
  const [game, setGame] = useState(null);
  const [guesses, setGuesses] = useState([]);
  const [currentGuess, setCurrentGuess] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [gameStatus, setGameStatus] = useState(''); // 'won', 'lost', 'playing'

  useEffect(() => {
    createNewGame();
  }, []);

  const createNewGame = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await api.createGame(6);
      setGame(response.game);
      setGuesses([]);
      setCurrentGuess('');
      setGameStatus('playing');
    } catch (err) {
      setError(`Failed to create game: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const submitGuess = async () => {
    if (!currentGuess || currentGuess.length !== 5) {
      setError('Please enter a 5-letter word');
      return;
    }

    if (!game) {
      setError('No active game');
      return;
    }

    setLoading(true);
    setError('');
    
    try {
      await api.makeGuess(game.id, currentGuess);
      
      // Refresh game state
      const gameState = await api.getGame(game.id);
      setGame(gameState.game);
      setGuesses(gameState.guesses || []);
      
      // Check game status
      if (gameState.game.is_won) {
        setGameStatus('won');
      } else if (gameState.game.is_completed) {
        setGameStatus('lost');
      }
      
      setCurrentGuess('');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const value = e.target.value.toUpperCase().replace(/[^A-Z]/g, '');
    if (value.length <= 5) {
      setCurrentGuess(value);
      setError('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      submitGuess();
    }
  };

  const renderGameBoard = () => {
    const rows = [];
    const maxGuesses = game?.max_guesses || 6;

    for (let i = 0; i < maxGuesses; i++) {
      const guess = guesses[i];
      const isCurrentRow = i === guesses.length && gameStatus === 'playing';
      
      rows.push(
        <div key={i} className="guess-row">
          {renderGuessRow(guess, isCurrentRow)}
        </div>
      );
    }

    return <div className="game-board">{rows}</div>;
  };

  const renderGuessRow = (guess, isCurrentRow) => {
    const tiles = [];
    
    for (let i = 0; i < 5; i++) {
      let letter = '';
      let status = '';
      
      if (guess && guess.result && guess.result[i]) {
        letter = guess.result[i].letter;
        status = guess.result[i].status;
      } else if (isCurrentRow && currentGuess[i]) {
        letter = currentGuess[i];
        status = 'filled';
      }
      
      tiles.push(
        <div key={i} className={`letter-tile ${status}`}>
          {letter}
        </div>
      );
    }
    
    return tiles;
  };

  const isGameOver = gameStatus === 'won' || gameStatus === 'lost';
  const canSubmit = currentGuess.length === 5 && !loading && !isGameOver;

  return (
    <div className="game-container">
      {loading && <div className="loading">Loading...</div>}
      
      {error && <div className="error-message">{error}</div>}
      
      {game && (
        <div className="game-info">
          <div>Game ID: {game.id}</div>
          <div>Guesses: {game.guess_count} / {game.max_guesses}</div>
          
          {gameStatus === 'won' && (
            <div className="game-status won">
              🎉 Congratulations! You won in {game.guess_count} guesses!
            </div>
          )}
          
          {gameStatus === 'lost' && (
            <div className="game-status lost">
              😞 Game over! The word was: {game.target_word}
            </div>
          )}
        </div>
      )}

      {game && renderGameBoard()}

      {!isGameOver && (
        <div className="input-section">
          <input
            type="text"
            value={currentGuess}
            onChange={handleInputChange}
            onKeyPress={handleKeyPress}
            placeholder="Enter your guess"
            className="guess-input"
            maxLength={5}
            disabled={loading}
          />
          <button
            onClick={submitGuess}
            disabled={!canSubmit}
            className="submit-btn"
          >
            Submit
          </button>
        </div>
      )}

      <div className="game-controls">
        <button onClick={createNewGame} className="new-game-btn" disabled={loading}>
          New Game
        </button>
      </div>
    </div>
  );
};

export default WordleGame;
