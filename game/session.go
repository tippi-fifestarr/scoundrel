package game

import (
	"errors"
	"sync"
	"time"
)

// SessionManager manages active game sessions
type SessionManager struct {
	sessions map[string]*GameSession
	mutex    sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*GameSession),
	}
}

// CreateSession creates a new game session
func (sm *SessionManager) CreateSession() string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session := NewGameSession()
	sm.sessions[session.GetID()] = session

	return session.GetID()
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(id string) (*GameSession, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessions[id]
	if !exists {
		return nil, errors.New("session not found")
	}

	return session, nil
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(id string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.sessions, id)
}

// PlayCard plays a card in the specified session
func (sm *SessionManager) PlayCard(sessionID string, cardIndex int) (*GameSession, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	err := session.PlayCard(cardIndex)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// PlayCardWithoutWeapon plays a card without using a weapon
func (sm *SessionManager) PlayCardWithoutWeapon(sessionID string, cardIndex int) (*GameSession, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	err := session.PlayCardWithoutWeapon(cardIndex)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// SkipRoom skips the current room in the specified session
func (sm *SessionManager) SkipRoom(sessionID string) (*GameSession, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	err := session.SkipRoom()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// CleanupSessions removes completed or stale sessions
func (sm *SessionManager) CleanupSessions(maxAge time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// In a real implementation, we would track creation time and
	// last access time for sessions and remove old ones

	// For now, just remove completed games
	for id, session := range sm.sessions {
		if session.IsGameOver() {
			delete(sm.sessions, id)
		}
	}
}

// ActiveSessionCount returns the number of active sessions
func (sm *SessionManager) ActiveSessionCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return len(sm.sessions)
}

// GetAllSessions returns all active sessions (for monitoring/debugging)
func (sm *SessionManager) GetAllSessions() []*GameSession {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sessions := make([]*GameSession, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}
