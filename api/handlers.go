package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tippi-fifestarr/scoundrel/game"
)

// Handler handles API requests for the game
type Handler struct {
	sessionManager *game.SessionManager
}

// NewHandler creates a new Handler
func NewHandler(sessionManager *game.SessionManager) *Handler {
	return &Handler{
		sessionManager: sessionManager,
	}
}

// CreateGameHandler creates a new game session
func (h *Handler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Create new game session
	sessionID := h.sessionManager.CreateSession()

	// Build response
	response := map[string]interface{}{
		"game_id": sessionID,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGameHandler returns the current state of a game
func (h *Handler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["id"]

	// Get game session
	session, err := h.sessionManager.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Game session not found", http.StatusNotFound)
		return
	}

	// Return game state
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.GetGameState())
}

// PlayCardHandler plays a card from the current room
func (h *Handler) PlayCardHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID and card index from URL
	vars := mux.Vars(r)
	sessionID := vars["id"]
	cardIndex, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid card index", http.StatusBadRequest)
		return
	}

	// Get game session
	session, err := h.sessionManager.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Game session not found", http.StatusNotFound)
		return
	}

	// Play the card
	err = session.PlayCard(cardIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated game state
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.GetGameState())
}

// PlayCardWithoutWeaponHandler plays a monster card without using an equipped weapon
func (h *Handler) PlayCardWithoutWeaponHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID and card index from URL
	vars := mux.Vars(r)
	sessionID := vars["id"]
	cardIndex, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid card index", http.StatusBadRequest)
		return
	}

	// Get game session
	session, err := h.sessionManager.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Game session not found", http.StatusNotFound)
		return
	}

	// Play the card without weapon
	err = session.PlayCardWithoutWeapon(cardIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated game state
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.GetGameState())
}

// SkipRoomHandler skips the current room
func (h *Handler) SkipRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["id"]

	// Get game session
	session, err := h.sessionManager.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Game session not found", http.StatusNotFound)
		return
	}

	// Skip the room
	err = session.SkipRoom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated game state
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.GetGameState())
}
