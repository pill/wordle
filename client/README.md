# Wordle Client

A React-based frontend for the Wordle game that connects to a Go backend API.

## Features

- Interactive Wordle game interface
- Real-time game state management
- Game statistics and history
- Modern, responsive design similar to the original Wordle
- Color-coded letter feedback (correct, present, absent)

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- Running Wordle Go server (backend)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm start
```

The app will open at `http://localhost:3000` and automatically proxy API requests to `http://localhost:8080`.

### Configuration

The client is configured to proxy requests to the backend server. If your backend runs on a different port, update the `proxy` field in `package.json`.

## API Integration

The client integrates with the following backend endpoints:

- `POST /api/games` - Create a new game
- `GET /api/games/{id}` - Get game state
- `POST /api/games/{id}` - Make a guess
- `GET /api/games` - Get recent games
- `GET /api/stats` - Get game statistics
- `DELETE /api/games/{id}` - Delete a game

## Components

- `App.js` - Main application component with navigation
- `WordleGame.js` - Core game component with game board and input
- `GameStats.js` - Statistics and recent games display
- `api.js` - API service for backend communication

## Styling

The app uses a dark theme similar to the New York Times Wordle with:
- Green tiles for correct letters in correct positions
- Yellow tiles for correct letters in wrong positions
- Gray tiles for letters not in the word
- Responsive design for mobile and desktop

## Available Scripts

- `npm start` - Run development server
- `npm build` - Build for production
- `npm test` - Run tests
- `npm eject` - Eject from Create React App

## Usage

1. The game automatically creates a new game when you load the page
2. Enter 5-letter words using the input field or by typing
3. Press Enter or click Submit to make a guess
4. View your game statistics in the Stats tab
5. Start a new game anytime with the "New Game" button
