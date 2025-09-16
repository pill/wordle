# Wordle

A Go implementation of Wordle game infrastructure.

## Requirements

- Game consists of
    - target word (5 letters)
    - guess number
    - guess history
    - 26 letter list (used, available, in word, in word in order)

## Prerequisites

- Go 1.21 or later

## Quick Start

```bash
# Clone the repository
git clone git@github.com:pill/wordle.git
cd wordle

# Run the word list demo
go run server/main.go server/wordlist.go

# Run tests
cd server && go test -v
```

## Server

The server component includes a robust word list management system implemented in Go:

- **Word List Management**: Load and validate Wordle words
- **Game Logic**: (planned) Start game, validate guesses, update board
- **Dictionary Operations**: Check word validity, get random words, filter by length

### Current Features

- ✅ Load 14,855 valid Wordle words
- ✅ Case-insensitive word validation
- ✅ Random word selection
- ✅ Word filtering by length
- ✅ Comprehensive test coverage
- ✅ Memory-efficient data structures

### Planned Features

- 🔄 Game state management
- 🔄 Turn-based gameplay API
- 🔄 Guess validation and scoring
- 🔄 Database integration for game history
- 🔄 Player statistics
- 🔄 REST API endpoints


### Display Features (Planned)
- letter used
- letter in word, but wrong spot
- letter in word, correct spot
- handle multiple letters
- show previous answers
- annotate answers (yellow, green)

## Development

### Running the Demo

```bash
# From project root
go run server/main.go server/wordlist.go
```

### Running Tests

```bash
# Test the word list functionality
cd server
go test -v

# Run specific tests
go test -v -run TestWordList
```

### Building

```bash
# Build the server
go build -o wordle-server server/*.go
```

## Word List Details

- **Source**: `server/valid-wordle-words.txt`
- **Count**: 14,855 five-letter words
- **Format**: One word per line, lowercase
- **Validation**: All standard Wordle-valid words included

## Contributing

1. Ensure Go 1.21+ is installed
2. Run tests before submitting PRs: `go test ./...`
3. Follow Go coding standards and conventions
4. Add tests for new functionality

## License

See LICENSE file for details.