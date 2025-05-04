# Scoundrel Card Game

A digital implementation of the Scoundrel card game, a solo roguelike dungeon crawler.

## Overview

Scoundrel is a challenging solo card game where you navigate through a dungeon made of cards, fighting monsters, finding weapons, and using potions. The game uses a modified standard deck of playing cards and features roguelike mechanics.

## Rules

### Setup
- A standard 52-card deck is used, with the red face cards and red aces removed (44 cards total)
- Player starts with 20 health points
- Goal: Make it through the entire dungeon without dying

### Card Types
- **Monsters** (Clubs & Spades): Deal damage equal to their value (2-14)
- **Weapons** (Diamonds): Reduce monster damage by their value
- **Health Potions** (Hearts): Restore health equal to their value (one effective potion per room)

### Gameplay
- Each room consists of 4 cards
- Player must play 3 cards in chosen order
- The remaining card starts the next room with 3 new cards
- Optional: Room may be skipped, but never twice in a row
- When using a weapon against monsters, the weapon can only be used against monsters with value less than or equal to the last monster it defeated

## Implementation

This project includes:

1. **Core Game Engine**: Written in Go with proper modeling of game mechanics
2. **Command-Line Interface**: For quick testing and gameplay
3. **RESTful API**: For interacting with the game programmatically
4. **Web Interface**: For a visual gameplay experience

## Project Structure

```
scoundrel/
├── cmd/                      # Application entry points
│   ├── api/                  # API server
│   │   └── main.go
│   └── cli/                  # Command-line interface
│       └── main.go
├── game/                     # Core game logic
│   ├── models.go             # Game models
│   ├── engine.go             # Game engine
│   └── session.go            # Session management
├── api/                      # API layer
│   ├── handlers.go           # API request handlers
│   └── server.go             # HTTP server setup
├── web/                      # Web frontend
│   ├── index.html            # Main HTML file
│   ├── js/                   # JavaScript files
│   │   ├── game.js           # Game logic
│   │   └── ui.js             # UI handling
│   └── css/                  # Stylesheets
│       └── styles.css        # Main styles
├── docs/                     # Documentation
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies
└── Makefile                  # Build commands
```

## Running the Game

### Prerequisites
- Go 1.16 or higher

### Command-Line Interface
To play the game in the terminal:

```bash
make run-cli
# or
go run cmd/cli/main.go
```

### API Server
To start the API server:

```bash
make run
# or
go run cmd/api/main.go
```

The server will start on http://localhost:8080

### Web Interface
To play the game with the web interface:

1. Start the API server as described above
2. Open http://localhost:8080 in your browser

## Development

### Testing
Run the tests:

```bash
make test
# or
go test ./...
```

### Building
Build the binaries:

```bash
make build
```

## Future Enhancements

1. AI Player implementation with different strategies
2. Blockchain integration for game state storage
3. Multiplayer functionality for competitive play
4. Enhanced animations and visual effects

## License

[MIT License](LICENSE)