import React, { useState } from 'react';
import WordleGame from './components/WordleGame';
import GameStats from './components/GameStats';
import './App.css';

function App() {
  const [currentView, setCurrentView] = useState('game'); // 'game' or 'stats'

  return (
    <div className="app">
      <header className="header">
        <h1 className="title">WORDLE</h1>
        <nav className="nav">
          <button 
            className={`nav-btn ${currentView === 'game' ? 'active' : ''}`}
            onClick={() => setCurrentView('game')}
          >
            Game
          </button>
          <button 
            className={`nav-btn ${currentView === 'stats' ? 'active' : ''}`}
            onClick={() => setCurrentView('stats')}
          >
            Stats
          </button>
        </nav>
      </header>

      <main className="main-content">
        {currentView === 'game' && <WordleGame />}
        {currentView === 'stats' && <GameStats />}
      </main>

      <footer className="footer">
        <p>Built with React - Connected to Go Wordle API</p>
      </footer>
    </div>
  );
}

export default App;
