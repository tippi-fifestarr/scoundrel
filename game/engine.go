package game

import (
	"errors"

	"github.com/google/uuid"
)

// GameState represents the current state of the game
type GameState int

const (
	GameStateInitial GameState = iota
	GameStateInProgress
	GameStateWon
	GameStateLost
)

// String returns a string representation of the game state
func (gs GameState) String() string {
	return [...]string{"Initial", "InProgress", "Won", "Lost"}[gs]
}

// GameSession represents an active game session
type GameSession struct {
	ID             string
	player         *Player
	deck           *Deck
	currentRoom    *Room
	playHistory    []*Card
	state          GameState
	lastCardPlayed *Card
}

// NewGameSession creates a new game session
func NewGameSession() *GameSession {
	id := uuid.New().String()
	player := NewPlayer(20) // Start with 20 health
	deck := NewDeck()
	deck.Shuffle()

	session := &GameSession{
		ID:          id,
		player:      player,
		deck:        deck,
		playHistory: make([]*Card, 0),
		state:       GameStateInitial,
	}

	// Create initial room
	session.CreateRoom()

	return session
}

// GetID returns the session ID
func (g *GameSession) GetID() string {
	return g.ID
}

// GetPlayer returns the player in this session
func (g *GameSession) GetPlayer() *Player {
	return g.player
}

// GetDeck returns the deck in this session
func (g *GameSession) GetDeck() *Deck {
	return g.deck
}

// GetCurrentRoom returns the current room
func (g *GameSession) GetCurrentRoom() *Room {
	return g.currentRoom
}

// GetState returns the current game state
func (g *GameSession) GetState() GameState {
	return g.state
}

// CreateRoom creates a new room in the dungeon
func (g *GameSession) CreateRoom() error {
	if g.state == GameStateWon || g.state == GameStateLost {
		return errors.New("game is already over")
	}

	// Create initial room or new room
	var cards []*Card
	var err error

	if g.currentRoom != nil && g.currentRoom.RemainingCard() != nil {
		// Use the remaining card from previous room
		remainingCard := g.currentRoom.RemainingCard()

		// Draw 3 more cards
		newCards, err := g.deck.Draw(3)
		if err != nil {
			// No more cards to draw, but we have a remaining card
			// This means the player has won by exhausting the deck
			g.state = GameStateWon
			return nil
		}

		cards = append([]*Card{remainingCard}, newCards...)
	} else {
		// Draw 4 cards for a new room
		cards, err = g.deck.Draw(4)
		if err != nil {
			// Can't draw enough cards, game is won
			g.state = GameStateWon
			return nil
		}
	}

	g.currentRoom = NewRoom(cards)
	g.state = GameStateInProgress
	g.player.SetUsedPotionThisRoom(false) // Reset potion usage for new room

	return nil
}

// PlayCard plays a card from the current room
func (g *GameSession) PlayCard(index int) error {
	if g.state != GameStateInProgress {
		return errors.New("game is not in progress")
	}

	// Play the card
	card, err := g.currentRoom.PlayCard(index)
	if err != nil {
		return err
	}

	// Save the last played card
	g.lastCardPlayed = card

	// Process card based on type
	switch card.Type() {
	case Monster:
		err = g.handleMonster(card)
	case Weapon:
		err = g.handleWeapon(card)
	case Potion:
		err = g.handlePotion(card)
	}

	if err != nil {
		return err
	}

	// Add to play history
	g.playHistory = append(g.playHistory, card)

	// Check if room is completed (3 cards played)
	if g.currentRoom.Completed() {
		// Set up next room
		err = g.CreateRoom()
		if err != nil {
			return err
		}
	}

	// Check if player is dead
	if g.player.Health() <= 0 {
		g.state = GameStateLost
	}

	return nil
}

// PlayCardWithoutWeapon plays a card from the current room without using a weapon
func (g *GameSession) PlayCardWithoutWeapon(index int) error {
	if g.state != GameStateInProgress {
		return errors.New("game is not in progress")
	}

	// Play the card
	card, err := g.currentRoom.PlayCard(index)
	if err != nil {
		return err
	}

	// Save the last played card
	g.lastCardPlayed = card

	// Process card based on type
	switch card.Type() {
	case Monster:
		err = g.handleMonsterWithoutWeapon(card)
	case Weapon:
		err = g.handleWeapon(card)
	case Potion:
		err = g.handlePotion(card)
	}

	if err != nil {
		return err
	}

	// Add to play history
	g.playHistory = append(g.playHistory, card)

	// Check if room is completed (3 cards played)
	if g.currentRoom.Completed() {
		// Set up next room
		err = g.CreateRoom()
		if err != nil {
			return err
		}
	}

	// Check if player is dead
	if g.player.Health() <= 0 {
		g.state = GameStateLost
	}

	return nil
}

// SkipRoom skips the current room
func (g *GameSession) SkipRoom() error {
	if g.state != GameStateInProgress {
		return errors.New("game is not in progress")
	}

	if g.deck.PrevRoomSkipped() {
		return errors.New("cannot skip two rooms in a row")
	}

	// Check if any cards have been played in the current room
	if len(g.currentRoom.playedCards) > 0 {
		return errors.New("cannot skip a room after playing cards")
	}

	// Add current room cards to bottom of deck
	g.deck.AddToBottom(g.currentRoom.AllCards())
	g.deck.SetPrevRoomSkipped(true)

	// Create a new room
	return g.CreateRoom()
}

// IsGameOver returns true if the game is over (won or lost)
func (g *GameSession) IsGameOver() bool {
	return g.state == GameStateWon || g.state == GameStateLost
}

// Handle monster card play
func (g *GameSession) handleMonster(card *Card) error {
	weapon := g.player.EquippedWeapon()

	// Check if we can use the weapon against this monster
	canUseWeapon := weapon != nil && g.player.CanUseWeaponAgainst(card)

	if canUseWeapon {
		// Calculate damage as monster value - weapon value
		damage := card.Value() - weapon.Value()
		if damage < 0 {
			damage = 0
		}

		// Apply damage to player
		err := g.player.ApplyDamage(damage)
		if err != nil {
			return err
		}

		// Add monster to defeated list for this weapon
		g.player.AddDefeatedMonster(card)
	} else {
		// Take full damage from monster
		err := g.player.ApplyDamage(card.Value())
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleMonsterWithoutWeapon forces fighting a monster without using a weapon
func (g *GameSession) handleMonsterWithoutWeapon(card *Card) error {
	// Take full damage from monster regardless of weapon
	err := g.player.ApplyDamage(card.Value())
	if err != nil {
		return err
	}

	return nil
}

// Handle weapon card play
func (g *GameSession) handleWeapon(card *Card) error {
	g.player.EquipWeapon(card)
	return nil
}

// Handle potion card play
func (g *GameSession) handlePotion(card *Card) error {
	// Only first potion in a room has effect
	if !g.player.UsedPotionThisRoom() {
		g.player.Heal(card.Value())
		g.player.SetUsedPotionThisRoom(true)
	}
	return nil
}

// GetGameState returns the current game state as a map for API responses
func (g *GameSession) GetGameState() map[string]interface{} {
	// Convert current room cards to a response format
	roomCards := make([]map[string]interface{}, 0)
	if g.currentRoom != nil {
		for i, card := range g.currentRoom.Cards() {
			roomCards = append(roomCards, map[string]interface{}{
				"index":   i,
				"suit":    int(card.Suit),
				"rank":    int(card.Rank),
				"value":   card.Value(),
				"type":    int(card.Type()),
				"display": card.String(),
			})
		}
	}

	// Convert equipped weapon to response format
	var equippedWeapon map[string]interface{}
	if g.player.EquippedWeapon() != nil {
		weapon := g.player.EquippedWeapon()
		equippedWeapon = map[string]interface{}{
			"suit":    int(weapon.Suit),
			"rank":    int(weapon.Rank),
			"value":   weapon.Value(),
			"display": weapon.String(),
		}
	}

	// Convert defeated monsters to response format
	defeatedMonsters := make([]map[string]interface{}, 0)
	for _, monster := range g.player.DefeatedMonsters() {
		defeatedMonsters = append(defeatedMonsters, map[string]interface{}{
			"suit":    int(monster.Suit),
			"rank":    int(monster.Rank),
			"value":   monster.Value(),
			"display": monster.String(),
		})
	}

	// Build and return complete game state
	return map[string]interface{}{
		"game_id": g.ID,
		"state":   g.state.String(),
		"player": map[string]interface{}{
			"health":            g.player.Health(),
			"max_health":        g.player.MaxHealth(),
			"equipped_weapon":   equippedWeapon,
			"defeated_monsters": defeatedMonsters,
			"used_potion":       g.player.UsedPotionThisRoom(),
		},
		"room": map[string]interface{}{
			"cards":     roomCards,
			"completed": g.currentRoom != nil && g.currentRoom.Completed(),
		},
		"deck": map[string]interface{}{
			"remaining_cards":       g.deck.Remaining(),
			"previous_room_skipped": g.deck.PrevRoomSkipped(),
		},
	}
}
