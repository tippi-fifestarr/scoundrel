package game

import (
	"testing"
)

func TestGameEngineFlow(t *testing.T) {
	session := NewGameSession()

	// Verify initial state
	if session.GetState() != GameStateInProgress {
		t.Errorf("Expected initial game state to be InProgress, got %v", session.GetState())
	}

	// Test playing cards one by one
	err := session.PlayCard(0) // Play first card
	if err != nil {
		t.Errorf("Error playing first card: %v", err)
	}

	// Check that the card was removed from the room
	if len(session.GetCurrentRoom().Cards()) != 3 {
		t.Errorf("Expected 3 cards left after playing 1, got %d", len(session.GetCurrentRoom().Cards()))
	}
}

func TestPotionMechanics(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Reduce player health so we can see the healing effect
	player.ApplyDamage(10) // Reduce to 10 health
	initialHealth := player.Health()

	// Set up test scenario with a controlled room
	potionCard1 := &Card{Suit: Hearts, Rank: Five}  // Potion with value 5
	potionCard2 := &Card{Suit: Hearts, Rank: Seven} // Potion with value 7
	monsterCard := &Card{Suit: Clubs, Rank: Ten}    // Monster with value 10
	weaponCard := &Card{Suit: Diamonds, Rank: Six}  // Weapon with value 6

	testRoom := NewRoom([]*Card{potionCard1, potionCard2, monsterCard, weaponCard})
	session.currentRoom = testRoom

	// Test first potion effect
	_ = session.handlePotion(potionCard1)

	if player.Health() != initialHealth+potionCard1.Value() {
		t.Errorf("Expected health to increase by %d, but got %d instead of %d",
			potionCard1.Value(), player.Health(), initialHealth+potionCard1.Value())
	}

	// Test second potion (should have no effect)
	healthAfterFirstPotion := player.Health()
	_ = session.handlePotion(potionCard2)

	if player.Health() != healthAfterFirstPotion {
		t.Errorf("Expected second potion to have no effect, but health changed from %d to %d",
			healthAfterFirstPotion, player.Health())
	}
}

func TestCombatMechanics(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Test monster without weapon
	monsterCard := &Card{Suit: Clubs, Rank: Eight} // Monster with value 8
	initialHealth := player.Health()
	_ = session.handleMonster(monsterCard)

	if player.Health() != initialHealth-monsterCard.Value() {
		t.Errorf("Expected to take full damage of %d, but took %d instead",
			monsterCard.Value(), initialHealth-player.Health())
	}

	// Test monster with weapon
	weaponCard := &Card{Suit: Diamonds, Rank: Five} // Weapon with value 5
	player.EquipWeapon(weaponCard)

	strongerMonster := &Card{Suit: Spades, Rank: Ten} // Monster with value 10
	healthBeforeSecondMonster := player.Health()
	_ = session.handleMonster(strongerMonster)

	expectedDamage := strongerMonster.Value() - weaponCard.Value()
	if player.Health() != healthBeforeSecondMonster-expectedDamage {
		t.Errorf("Expected to take reduced damage of %d, but took %d instead",
			expectedDamage, healthBeforeSecondMonster-player.Health())
	}

	// Test weapon restriction after first monster defeated
	weakerMonster := &Card{Suit: Clubs, Rank: Nine} // Monster with value 9
	_ = session.handleMonster(weakerMonster)

	// Test weapon restriction after second monster defeated
	evenWeakerMonster := &Card{Suit: Spades, Rank: Eight} // Monster with value 8
	healthBeforeLastMonster := player.Health()
	_ = session.handleMonster(evenWeakerMonster)

	// Should be able to use weapon because 8 < 9
	expectedDamage = evenWeakerMonster.Value() - weaponCard.Value()
	if player.Health() != healthBeforeLastMonster-expectedDamage {
		t.Errorf("Expected to take reduced damage of %d against weaker monster, but took %d instead",
			expectedDamage, healthBeforeLastMonster-player.Health())
	}

	// Test weapon restriction violation with stronger monster
	evenStrongerMonster := &Card{Suit: Clubs, Rank: Queen} // Monster with value 12
	player.health = 20                                     // Reset health for this test
	healthBeforeStrongMonster := player.Health()
	_ = session.handleMonster(evenStrongerMonster)

	// Should NOT be able to use weapon because Queen value > all defeated monsters
	if player.Health() != healthBeforeStrongMonster-evenStrongerMonster.Value() {
		t.Errorf("Expected to take full damage of %d against stronger monster, but took %d instead",
			evenStrongerMonster.Value(), healthBeforeStrongMonster-player.Health())
	}
}

func TestRoomSkipping(t *testing.T) {
	session := NewGameSession()

	// Should be able to skip first room
	err := session.SkipRoom()
	if err != nil {
		t.Errorf("Expected to be able to skip first room, but got error: %v", err)
	}

	// Should not be able to skip second room (no skipping two rooms in a row)
	err = session.SkipRoom()
	if err == nil {
		t.Errorf("Expected error when skipping two rooms in a row, but got none")
	}
}

func TestRoomProgression(t *testing.T) {
	// Create a custom session with a controlled deck for testing
	session := NewGameSession()

	// Create a mock deck with sufficient cards
	mockDeck := &Deck{
		cards: make([]*Card, 20), // Plenty of cards for testing
	}

	// Fill deck with identifiable cards
	for i := 0; i < 20; i++ {
		rank := Rank((i % 13) + 2) // 2-14 (2-Ace)
		suit := Suit(i % 4)        // 0-3 (Clubs, Diamonds, Hearts, Spades)
		mockDeck.cards[i] = &Card{Suit: suit, Rank: rank}
	}

	// Replace session's deck
	session.deck = mockDeck

	// Create a test room
	room := NewRoom([]*Card{
		{Suit: Hearts, Rank: Two},     // Potion
		{Suit: Diamonds, Rank: Three}, // Weapon
		{Suit: Clubs, Rank: Four},     // Monster
		{Suit: Spades, Rank: Five},    // Monster
	})
	session.currentRoom = room

	// Test playing cards one by one
	err := session.PlayCard(0) // Play potion
	if err != nil {
		t.Errorf("Error playing first card: %v", err)
	}

	if len(session.currentRoom.Cards()) != 3 {
		t.Errorf("Expected 3 cards left after playing 1, got %d", len(session.currentRoom.Cards()))
	}

	err = session.PlayCard(0) // Play weapon
	if err != nil {
		t.Errorf("Error playing second card: %v", err)
	}

	if len(session.currentRoom.Cards()) != 2 {
		t.Errorf("Expected 2 cards left after playing 2, got %d", len(session.currentRoom.Cards()))
	}

	err = session.PlayCard(0) // Play monster
	if err != nil {
		t.Errorf("Error playing third card: %v", err)
	}

	// After playing 3 cards, we should have created a new room with the 4th card + 3 new cards
	if len(session.currentRoom.Cards()) != 4 {
		t.Errorf("Expected a new room with 4 cards after completing previous room, got %d",
			len(session.currentRoom.Cards()))
	}
}
