# Scoundrel Card Game - Implementation Tutorial

This tutorial guides you through the implementation of the Scoundrel card game, explaining the core concepts, architecture, and code structure. This document is intended for developers who want to understand how the game is built or contribute to its development.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Game Rules](#game-rules)
3. [Project Structure](#project-structure)
4. [Core Data Models](#core-data-models)
5. [Game Engine Implementation](#game-engine-implementation)
6. [Session Management](#session-management)
7. [API Layer](#api-layer)
8. [CLI Implementation](#cli-implementation)
9. [Building and Running](#building-and-running)

## Project Overview

Scoundrel is a single-player roguelike dungeon crawler card game implemented as a Go application with a REST API. The implementation includes:

- Core game logic
- Session management for multiple concurrent games
- HTTP API for interfacing with the game
- Command-line interface for direct play

We've used a simplified architecture to accelerate development, focusing on core functionality first.

## Game Rules

Scoundrel is played with a modified standard deck of cards:

- **Setup**: Start with 20 health points and a 42-card deck (standard 52-card deck with red face cards and aces removed)
- **Room Creation**: 4 cards form a room
- **Gameplay**: Play 3 cards from each room in your chosen order
- **Progression**: The remaining card starts the next room with 3 new cards
- **Room Skipping**: Optional, but never twice in a row

**Card Types**:
- **Monsters** (Clubs & Spades): Deal damage equal to their value (2-14)
- **Weapons** (Diamonds): Reduce monster damage by their value
- **Health Potions** (Hearts): Restore health equal to their value (one effective potion per room)

**Combat Rules**:
- Fighting without a weapon: Take full monster damage
- Fighting with a weapon: Take (monster value - weapon value) damage
- Weapon restriction: After defeating a monster, a weapon can only be used against monsters with lower values

**Win/Loss Conditions**:
- Win by exhausting the dungeon deck without dying
- Lose if health reaches 0

## Project Structure

We've organized the project into a simplified structure:

```
scoundrel/
├── cmd/                      # Application entry points
│   ├── api/                  # API server
│   │   └── main.go
│   └── cli/                  # CLI for testing
│       └── main.go
├── game/                     # Core game logic
│   ├── models.go             # All game models
│   ├── engine.go             # Game engine and rules
│   └── session.go            # Game session management
├── api/                      # API layer
│   ├── handlers.go           # Request handlers
│   └── server.go             # HTTP server setup
├── go.mod                    # Go module definition
├── go.sum                    # Go dependency checksums
├── Makefile                  # Build commands
└── README.md                 # Project documentation
```

## Core Data Models

The core data models are implemented in `game/models.go` and include:

### Card Model

Cards in Scoundrel are represented by suit and rank:

```go
type Suit int
const (
    Clubs Suit = iota
    Diamonds
    Hearts
    Spades
)

type Rank int
const (
    Two Rank = iota + 2
    Three
    // ...
    Ace = 14
)

type Card struct {
    Suit  Suit
    Rank  Rank
}
```

Each card has a value (its rank value) and a type based on its suit:
- Clubs/Spades → Monsters
- Diamonds → Weapons
- Hearts → Potions

### Player Model

The Player model tracks health, equipped weapon, and game state:

```go
type Player struct {
    health            int
    maxHealth         int
    equippedWeapon    *Card
    defeatedMonsters  []*Card
    usedPotionThisRoom bool
}
```

Key functionalities include:
- Health management (apply damage, heal)
- Weapon equipment
- Tracking monsters defeated with the current weapon
- Verifying weapon usage restrictions
- Tracking potion usage per room

### Room Model

The Room represents a set of cards the player must navigate:

```go
type Room struct {
    cards        []*Card
    playedCards  []*Card
}
```

It handles:
- Playing cards from the room
- Tracking which cards have been played
- Identifying the remaining card for the next room
- Determining when a room is completed

### Deck Model

The Deck manages the cards that form the dungeon:

```go
type Deck struct {
    cards           []*Card
    prevRoomSkipped bool
}
```

It provides functionality for:
- Creating a standard deck (minus red face cards and aces)
- Shuffling cards
- Drawing cards
- Adding cards to the bottom of the deck
- Tracking if the previous room was skipped

## Game Engine Implementation

The game engine (`game/engine.go`) implements the core game logic:

```go
type GameSession struct {
    ID           string
    player       *Player
    deck         *Deck
    currentRoom  *Room
    playHistory  []*Card
    state        GameState
}
```

Key features include:

1. **Game States**:
   - Initial
   - InProgress
   - Won
   - Lost

2. **Room Creation**:
   - Uses the remaining card from the previous room
   - Draws 3 new cards
   - Creates a fresh room when starting

3. **Card Play Logic**:
   - Handles different card types (monsters, weapons, potions)
   - Applies combat rules
   - Manages room progression

4. **Combat Resolution**:
   - Calculates damage based on monster value and weapon value
   - Applies weapon usage restrictions
   - Updates player health

5. **Game Flow Control**:
   - Tracking game state
   - Detecting win/loss conditions
   - Managing room skipping

## Session Management

Session management (`game/session.go`) provides a way to handle multiple concurrent game sessions:

```go
type SessionManager struct {
    sessions map[string]*GameSession
    mutex    sync.RWMutex
}
```

This manager:
- Creates new game sessions with unique IDs
- Retrieves sessions by ID
- Provides thread-safe access to session data
- Handles session cleanup
- Exposes methods for game actions (play card, skip room)

The implementation uses a mutex to ensure thread safety when accessing the session map, making it suitable for use in a concurrent environment like a web server.

## API Layer

The API layer consists of two main components:

### Handlers (`api/handlers.go`)

Handlers process HTTP requests and interact with the game sessions:

```go
type Handler struct {
    sessionManager *game.SessionManager
}
```

Implemented endpoints:
- `POST /api/games` - Create a new game
- `GET /api/games/{id}` - Get game state
- `POST /api/games/{id}/play/{index}` - Play a card
- `POST /api/games/{id}/skip` - Skip a room

Each handler validates the request, performs the requested action using the session manager, and returns an appropriate JSON response.

### Server (`api/server.go`)

The server sets up routes and manages the HTTP server:

```go
type Server struct {
    router         *mux.Router
    handler        *Handler
    sessionManager *game.SessionManager
}
```

It configures:
- Route definitions using gorilla/mux
- Middleware for logging and CORS
- HTTP server settings

## CLI Implementation

The CLI (`cmd/cli/main.go`) provides a text-based interface for playing the game:

```go
func main() {
    session := game.NewGameSession()
    // Game loop
    for !session.IsGameOver() {
        displayGameState(session)
        action := getPlayerAction()
        executeAction(action, session)
    }
}
```

It implements:
1. Game state display (health, equipped weapon, current room)
2. Player action input
3. Action execution (play card, skip room)
4. Game loop control

The CLI is useful for testing the game logic directly and provides an alternative way to play the game without using the API.

## Building and Running

### Setting Up Go

We started by ensuring Go was installed on the system:

```bash
brew install go  # On macOS
```

### Initializing the Project

We initialized a Go module for the project:

```bash
go mod init github.com/tippi-fifestarr/scoundrel
```

### Adding Dependencies

We added the required external dependencies:

```bash
go get github.com/gorilla/mux    # For HTTP routing
go get github.com/google/uuid    # For generating unique IDs
```

### Setting Up the Project Structure

We created the necessary directories:

```bash
mkdir -p cmd/api cmd/cli game api
```

### Implementing the Core Components

We implemented the components in this order:
1. Core models (`game/models.go`)
2. Game engine (`game/engine.go`)
3. Session management (`game/session.go`)
4. API handlers (`api/handlers.go`)
5. API server (`api/server.go`)
6. API entrypoint (`cmd/api/main.go`)
7. CLI interface (`cmd/cli/main.go`)
8. Makefile for build commands

### Updating Dependencies

Finally, we make sure all dependencies are properly recorded in go.mod and go.sum:

```bash
go mod tidy
```

## Next Steps

After completing the initial implementation, you can:

1. Build and run the CLI to test the game:
   ```bash
   make run-cli
   ```

2. Run the API server to interact via HTTP:
   ```bash
   make run-api
   ```

3. Implement a web frontend to provide a graphical interface
4. Add automated tests for game logic
5. Implement AI players
6. Add blockchain integration

## Conclusion

This tutorial has covered the implementation of the Scoundrel card game in Go, explaining the core concepts, code structure, and development process. The game implements all the rules of Scoundrel while providing both API and CLI interfaces.

The simplified architecture allows for rapid development while maintaining flexibility for future expansion. By organizing the code into logical components, we've created a maintainable and extendable codebase.