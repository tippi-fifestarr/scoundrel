# Scoundrel Game Backend

This repository contains a Go implementation of the Scoundrel card game, designed to support concurrent gameplay with both human and AI players.

## Core Components

The backend consists of three main packages:

### `/game` Directory Structure

```
/game
  ├── models.go    # Core game entities and data structures
  ├── engine.go    # Game rules and session logic
  └── session.go   # Concurrent session management
```

## Concurrency Architecture

This backend implements a thread-safe concurrent model that allows multiple game sessions to run simultaneously. Key features include:

- **Thread Isolation**: Each game session runs independently
- **Thread-Safe Access**: All shared resources are protected by appropriate locks
- **Read/Write Optimizations**: Uses RWMutex to allow concurrent reads
- **Safe State Updates**: Game state changes are properly synchronized

## File Descriptions

### `models.go`

Contains the core game data structures and entities:

- **Card**: Represents a playing card with suit and rank
- **Player**: Tracks player state including health, weapons, and potion usage
- **Room**: Represents a dungeon room with cards
- **Deck**: Manages the collection of cards for gameplay

These models form the foundation of the game but do not handle concurrency directly.

### `engine.go`

Implements the core game logic and state management:

- **GameSession**: Represents a single game instance
- **GameState**: State machine tracking game progress (Initial, InProgress, Won, Lost)
- **Game Actions**: Methods for playing cards, creating rooms, handling monster combat

The engine is designed to be used within a single thread context and relies on the session manager for thread safety.

### `session.go`

Manages concurrent access to game sessions:

- **SessionManager**: Thread-safe manager for multiple concurrent game sessions
- **Concurrency Control**: Uses sync.RWMutex to control access to shared session data
- **Session Lifecycle**: Methods for creating, retrieving, and removing sessions

## Concurrency Design Patterns

### Reader-Writer Lock Pattern

```go
// Read operation example
func (sm *SessionManager) GetSession(id string) (*GameSession, error) {
    sm.mutex.RLock()       // Acquire read lock
    defer sm.mutex.RUnlock() // Ensure unlock happens
    
    session, exists := sm.sessions[id]
    if !exists {
        return nil, errors.New("session not found")
    }
    
    return session, nil
}

// Write operation example
func (sm *SessionManager) CreateSession() *GameSession {
    sm.mutex.Lock()       // Acquire exclusive lock
    defer sm.mutex.Unlock() // Ensure unlock happens
    
    session := NewGameSession()
    sm.sessions[session.ID] = session
    
    return session
}
```

### Session Isolation

Each GameSession operates independently, allowing for concurrent gameplay without contention. The SessionManager only synchronizes access to the sessions map, not to the individual game state of each session.

## Performance Considerations

1. **Lock Granularity**: Locks are applied at the session manager level, not individual sessions
2. **Read Optimization**: Read locks allow multiple concurrent reads
3. **Lock Duration**: Locks are held only for the minimum time necessary
4. **Cleanup**: Stale sessions are periodically removed to prevent memory leaks

## Future Concurrency Extensions Ideas

To support AI players:

1. **AI Orchestrator**: Will manage concurrent AI game simulations
2. **Batch Processing**: Will allow multiple AI games to run in parallel
3. **Training Pipeline**: Will collect game results for model improvement

## Thread Safety Guide

When extending the system:

1. Access all sessions through the SessionManager, never directly
2. Keep lock durations as short as possible
3. Use read locks for operations that don't modify the sessions map
4. Use write locks for operations that modify the sessions map
5. Consider using atomic operations for simple counters and flags