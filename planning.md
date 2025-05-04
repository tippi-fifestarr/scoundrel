# Scoundrel Game Implementation Plan

## Project Overview

Scoundrel is a solo roguelike dungeon crawler card game that we're implementing as a full-stack application with future AI and blockchain integration. This implementation plan focuses on a simplified approach to get the core functionality working quickly.

## Simplified Project Structure

```
scoundrel/
├── cmd/                      # Application entry points
│   ├── api/                  # API server
│   │   └── main.go
│   └── cli/                  # Simple CLI for testing
│       └── main.go
├── game/                     # Core game logic
│   ├── models.go             # All game models in one file
│   ├── engine.go             # Game engine and rules
│   └── session.go            # Game session management
├── api/                      # Simple API layer
│   ├── handlers.go           # API request handlers
│   └── server.go             # HTTP server setup
├── web/                      # Frontend (upcoming)
│   ├── index.html
│   ├── js/
│   │   ├── game.js
│   │   └── ui.js
│   └── css/
│       └── styles.css
├── docs/                     # Documentation
│   ├── tutorial.md           # Developer tutorial
│   └── rules.md              # Game rules and test cases
├── go.mod
├── go.sum
├── Makefile                  # Simple build commands
└── README.md
```

## Implementation Phases

### Phase 1: Core Game Logic (2 weeks) ✅

**Week 1: Base Implementation** ✅
- Implement card, deck, and player models in a single file
- Create basic game mechanics and rules
- Implement a simple in-memory session manager

**Week 2: Game Flow & Testing** ✅
- Complete room progression logic
- Add combat resolution
- Create basic CLI for manual testing
- Add core unit tests

### Phase 2: Simple API (2 weeks) ✅

**Week 3: HTTP API** ✅
- Create basic HTTP handlers for game actions
- Implement JSON serialization of game state
- Add simple middleware for logging and error handling

**Week 4: Integration & Testing** ✅
- Connect API to game engine
- Add basic authentication
- Create integration tests
- Build simple test client

### Phase 3: Web Frontend (2 weeks) 🔄

**Week 5: Basic UI**
- Create HTML structure and CSS styling
- Implement card visualization
- Design game board layout
- Add health and status displays

**Week 6: Frontend Logic**
- Connect frontend to API
- Implement game flow in JavaScript
- Add animations and visual feedback
- Handle player interactions

## Development Timeline

```
May 5-6:    Project setup and initial structure          ✅
May 7-9:    Implement core models (Card, Deck, Player)   ✅
May 10-13:  Implement game engine logic                  ✅
May 14-16:  Basic HTTP API implementation                ✅
May 17-18:  Session management                           ✅
May 19-21:  Unit tests                                   ✅
May 22-23:  Manual testing and bug fixes                 ✅
May 24-28:  Web frontend - Basic UI                      🔄
May 29-31:  Web frontend - Game logic                    🔄
```

## Key Implementation Components

### 1. Card and Deck Models ✅

The card and deck models are the foundation of the game:

- Cards have suit, rank, and value
- Cards are categorized as monsters (clubs/spades), weapons (diamonds), or potions (hearts)
- Deck contains 44 cards (52 minus 8 red face cards and aces)
- Deck supports drawing cards and shuffling

### 2. Player Model ✅

The player model tracks:

- Current and maximum health (starts at 20)
- Currently equipped weapon
- Monsters defeated by the current weapon
- Potion usage flag for the current room

### 3. Room Mechanics ✅

Room implementation handles:

- 4 cards per room
- Playing 3 cards in player-chosen order
- Tracking the remaining card for the next room
- Handling room completion

### 4. Combat System ✅

The combat system implements:

- Monster damage calculation
- Weapon usage against monsters (with option to fight barehanded)
- Weapon restrictions (can only be used against monsters with lower value than previously defeated monsters)
- Health tracking and game-over detection

### 5. Simple REST API ✅

The API provides endpoints for:

- Creating a new game
- Getting current game state
- Playing a card
- Optional: Skipping a room

### 6. Web Frontend 🔄

The web frontend will include:

- Card visualization with appropriate styling
- Game board layout showing current room
- Player status display (health, weapon, etc.)
- Interactive controls for playing cards and making decisions
- Visual feedback for game events

## Testing Strategy

Our testing approach focuses on core functionality first:

1. Unit tests for game models and rules ✅
2. API tests for endpoints ✅
3. Manual testing through CLI ✅
4. End-to-end testing with frontend (upcoming)

## Tradeoffs Made for Simplicity

1. Consolidated directory structure with fewer packages ✅
2. Combined related code into fewer, larger files ✅
3. Simple REST API without WebSockets initially ✅
4. Focus on core unit tests with more manual testing ✅
5. Delayed implementation of advanced features ✅

## Next Steps After Frontend Implementation

1. AI player implementation
2. Performance optimization
3. Blockchain integration
4. Advanced game features and variations

## Development Tools

### Makefile

```makefile
.PHONY: run test build

run:
	go run cmd/api/main.go

test:
	go test ./...

build:
	go build -o bin/scoundrel-api cmd/api/main.go
```

## Future Expansion Areas

1. **Frontend Development** 🔄
   - React-based game interface
   - Card animations and game visualization
   - Game state management

2. **AI Integration**
   - Rule-based AI players
   - Learning-based AI strategies
   - Game analysis tools

3. **Blockchain Features**
   - NFT integration for cards
   - Game record storage
   - On-chain leader mechanics

## Implementation Progress

### Phase 1: Core Game Logic ✅
- [x] Implement card, deck, and player models in a single file
- [x] Create basic game mechanics and rules
- [x] Implement a simple in-memory session manager
- [x] Complete room progression logic
- [x] Add combat resolution
- [x] Create basic CLI for manual testing
- [x] Add core unit tests

### Phase 2: Simple API ✅
- [x] Create basic HTTP handlers for game actions
- [x] Implement JSON serialization of game state
- [x] Add simple middleware for logging and error handling
- [x] Connect API to game engine
- [x] Create integration tests

### Phase 3: Web Frontend 🔄
- [ ] Create HTML structure and CSS styling
- [ ] Implement card visualization
- [ ] Design game board layout
- [ ] Add health and status displays
- [ ] Connect frontend to API
- [ ] Implement game flow in JavaScript

### Documentation ✅
- [x] Create comprehensive README
- [x] Implement project planning document
- [x] Create developer tutorial
- [x] Add detailed rules documentation
- [x] Add code comments and documentation

## Challenges Addressed and Improvements Made

1. **Weapon Usage Choice Implementation** ✅
   - Added ability for players to choose between using weapon or fighting barehanded
   - Updated CLI to prompt players when facing monsters
   - Enhanced game engine to support both combat approaches

2. **Weapon Restriction Rule Correction** ✅
   - Updated the rule to check against the most recently defeated monster only
   - Fixed test cases to verify correct behavior
   - Improved documentation to clarify the "accordion" style rule

3. **Win Condition Fix** ✅
   - Fixed game state transition to properly detect when the player has won
   - Enhanced the CreateRoom method to set game state to Won when the deck is exhausted
   - Added test cases to verify correct win detection

4. **Health Cap Implementation** ✅
   - Updated the Heal method to enforce the maximum health of 20
   - Added clear documentation about the health cap rule
   - Fixed tests to account for health cap

5. **Room Skipping Enhancement** ✅
   - Added validation to prevent skipping after cards have been played
   - Updated CLI to only show skip option when applicable
   - Added error messages to clarify skip restrictions

## Frontend Implementation Plan

For the web frontend, we'll take a progressive approach:

1. **Basic HTML/CSS Layout**
   - Create the game board structure
   - Style cards for each type (monster, weapon, potion)
   - Design health and status indicators

2. **Card Visualization**
   - Render cards with appropriate suit and rank
   - Add visual differentiation for card types
   - Create a discard pile area

3. **Game State Management**
   - Connect to the API to retrieve game state
   - Maintain local state for smooth interaction
   - Handle game state transitions

4. **User Interactions**
   - Implement card selection
   - Create UI for weapon choice when fighting monsters
   - Add room skipping controls

5. **Visual Feedback**
   - Add animations for card play
   - Create visual effects for combat
   - Provide feedback for healing and weapon equipping

## Questions for Further Development

1. How complex should the frontend be initially? Basic HTML/JS or a full React application?
2. Should we implement real-time updates via WebSockets after the basic frontend is working?
3. Are there any additional game mechanics or rules you'd like to implement beyond the core rules?
4. How should we approach the AI implementation after the frontend is complete?