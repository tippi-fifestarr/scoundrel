package game

import (
	"testing"
)

// TestOptionalWeaponUse verifies that players can choose not to use their equipped weapon
func TestOptionalWeaponUse(t *testing.T) {
	// This rule isn't currently implemented in our game engine
	// We should modify the engine to support this feature

	t.Skip("Not implemented yet: Optional weapon usage needs to be added to the game engine")

	/*
		The implementation would look something like this:

		session := NewGameSession()
		player := session.GetPlayer()

		// Equip a weapon
		weaponCard := &Card{Suit: Diamonds, Rank: Eight} // Weapon with value 8
		player.EquipWeapon(weaponCard)

		// Face a monster
		monsterCard := &Card{Suit: Clubs, Rank: Ten} // Monster with value 10
		initialHealth := player.Health()

		// Choose to fight without weapon despite having one
		session.handleMonsterWithoutWeapon(monsterCard)

		// Should take full damage, not reduced damage
		if player.Health() != initialHealth-monsterCard.Value() {
			t.Errorf("Expected to take full damage of %d when choosing not to use weapon, but took %d instead",
			         monsterCard.Value(), initialHealth-player.Health())
		}
	*/
}

// TestWeaponAccordionRule tests the rule that weapons can only be used against
// progressively weaker monsters
func TestWeaponAccordionRule(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Create test cards
	weaponCard := &Card{Suit: Diamonds, Rank: Ten}      // Weapon with value 10
	strongMonster := &Card{Suit: Clubs, Rank: Queen}    // Monster with value 12
	weakerMonster := &Card{Suit: Spades, Rank: Nine}    // Monster with value 9
	evenWeakerMonster := &Card{Suit: Clubs, Rank: Five} // Monster with value 5
	strongerMonster := &Card{Suit: Spades, Rank: King}  // Monster with value 13

	// Equip weapon
	player.EquipWeapon(weaponCard)

	// First monster - weapon works (takes reduced damage)
	initialHealth := player.Health()
	session.handleMonster(strongMonster)
	expectedDamage := strongMonster.Value() - weaponCard.Value()
	if player.Health() != initialHealth-expectedDamage {
		t.Errorf("Expected damage to be reduced by weapon (damage: %d), but took %d damage",
			expectedDamage, initialHealth-player.Health())
	}

	// Second monster - should work because 9 < 12
	initialHealth = player.Health()
	session.handleMonster(weakerMonster)
	expectedDamage = weakerMonster.Value() - weaponCard.Value()
	if expectedDamage < 0 {
		expectedDamage = 0
	}
	if player.Health() != initialHealth-expectedDamage {
		t.Errorf("Expected to take reduced damage of %d against weaker monster, but took %d instead",
			expectedDamage, initialHealth-player.Health())
	}

	// Third monster - should work because 5 < 9
	initialHealth = player.Health()
	session.handleMonster(evenWeakerMonster)
	expectedDamage = evenWeakerMonster.Value() - weaponCard.Value()
	if expectedDamage < 0 {
		expectedDamage = 0
	}
	if player.Health() != initialHealth-expectedDamage {
		t.Errorf("Expected to take reduced damage of %d against even weaker monster, but took %d instead",
			expectedDamage, initialHealth-player.Health())
	}

	// Fourth monster - should NOT work because 13 > 5 (the last defeated monster)
	initialHealth = player.Health()
	session.handleMonster(strongerMonster)
	if player.Health() != initialHealth-strongerMonster.Value() {
		t.Errorf("Expected to take full damage of %d against stronger monster, but took %d instead",
			strongerMonster.Value(), initialHealth-player.Health())
	}
}

// TestWeaponReplacement verifies that equipping a new weapon resets the defeated monsters history
func TestWeaponReplacement(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Create test cards
	weaponCard1 := &Card{Suit: Diamonds, Rank: Five}  // Weapon with value 5
	weaponCard2 := &Card{Suit: Diamonds, Rank: Eight} // Weapon with value 8
	monsterCard := &Card{Suit: Clubs, Rank: Ten}      // Monster with value 10
	strongMonster := &Card{Suit: Spades, Rank: King}  // Monster with value 13

	// Equip first weapon and defeat a monster
	player.EquipWeapon(weaponCard1)
	session.handleMonster(monsterCard)

	// Check that the monster was recorded
	if len(player.DefeatedMonsters()) != 1 {
		t.Errorf("Expected 1 defeated monster, but got %d", len(player.DefeatedMonsters()))
	}

	// Try to defeat a stronger monster (should fail = full damage)
	initialHealth := player.Health()
	session.handleMonster(strongMonster)
	if player.Health() != initialHealth-strongMonster.Value() {
		t.Errorf("Expected to take full damage (%d) against stronger monster, but took %d instead",
			strongMonster.Value(), initialHealth-player.Health())
	}

	// Equip new weapon, history should reset
	player.EquipWeapon(weaponCard2)
	if len(player.DefeatedMonsters()) != 0 {
		t.Errorf("Expected defeated monsters to be reset after equipping new weapon, but got %d monsters",
			len(player.DefeatedMonsters()))
	}

	// Now should be able to defeat the stronger monster with reduced damage
	initialHealth = player.Health()
	session.handleMonster(strongMonster)
	expectedDamage := strongMonster.Value() - weaponCard2.Value()
	if player.Health() != initialHealth-expectedDamage {
		t.Errorf("Expected to take reduced damage of %d with new weapon, but took %d instead",
			expectedDamage, initialHealth-player.Health())
	}
}

// TestPotionReset verifies that potion usage flag resets between rooms
func TestPotionReset(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Reduce health for testing healing
	player.ApplyDamage(10) // Health now 10
	initialHealth := player.Health()

	// Create a room with multiple potions
	potionCard1 := &Card{Suit: Hearts, Rank: Five}  // Potion with value 5
	potionCard2 := &Card{Suit: Hearts, Rank: Seven} // Potion with value 7
	monsterCard1 := &Card{Suit: Clubs, Rank: Three} // Monster with value 3
	monsterCard2 := &Card{Suit: Spades, Rank: Four} // Monster with value 4

	// Create custom room for testing
	customRoom := NewRoom([]*Card{potionCard1, potionCard2, monsterCard1, monsterCard2})
	session.currentRoom = customRoom

	// Use first potion - should heal
	session.handlePotion(potionCard1)
	expectedHealth := initialHealth + potionCard1.Value()
	if player.Health() != expectedHealth {
		t.Errorf("Expected health to be %d after first potion, but got %d",
			expectedHealth, player.Health())
	}

	// Use second potion - should have no effect
	session.handlePotion(potionCard2)
	if player.Health() != expectedHealth {
		t.Errorf("Expected health to remain %d after second potion, but got %d",
			expectedHealth, player.Health())
	}

	// Complete room by playing all cards
	// (This would normally be done via PlayCard, but we're simulating it)
	player.SetUsedPotionThisRoom(false) // Simulate entering a new room

	// In a new room, potions should work again
	initialHealth = player.Health()
	session.handlePotion(potionCard2)

	// Remember, health cannot exceed 20 according to official rules
	expectedHealth = initialHealth + potionCard2.Value()
	if expectedHealth > 20 {
		expectedHealth = 20 // Cap at max health
	}

	if player.Health() != expectedHealth {
		t.Errorf("Expected health to be %d after potion in new room, but got %d",
			expectedHealth, player.Health())
	}
}

// TestWinCondition verifies that the game is won when the deck is exhausted
func TestWinCondition(t *testing.T) {
	session := NewGameSession()

	// Create a small deck to easily test win condition
	mockDeck := &Deck{
		cards: make([]*Card, 0), // Empty deck
	}
	session.deck = mockDeck

	// Set up room with exactly 4 cards
	weaponCard := &Card{Suit: Diamonds, Rank: Five} // Weapon
	potionCard := &Card{Suit: Hearts, Rank: Seven}  // Potion
	monsterCard1 := &Card{Suit: Clubs, Rank: Three} // Monster
	monsterCard2 := &Card{Suit: Spades, Rank: Four} // Monster

	// We need to make sure the room is properly created
	room := NewRoom([]*Card{weaponCard, potionCard, monsterCard1, monsterCard2})
	session.currentRoom = room

	// Since we're testing win condition, our deck is empty
	// So when we try to create a new room after playing 3 cards,
	// there will be no cards to draw and the game should be won

	// Play 3 cards manually
	session.PlayCard(0) // First card
	session.PlayCard(0) // Second card
	session.PlayCard(0) // Third card - this should trigger the win condition

	// The game should change state to Won when all cards are played
	if session.GetState() != GameStateWon {
		t.Logf("Current state: %v", session.GetState())
		t.Errorf("Expected game state to be Won after exhausting deck, but got %v", session.GetState())
	}
}

// TestLossCondition verifies that the game is lost when health reaches 0
func TestLossCondition(t *testing.T) {
	session := NewGameSession()
	player := session.GetPlayer()

	// Reduce health to just above 0
	player.ApplyDamage(19) // Health now 1

	// Create a monster that will kill the player
	monsterCard := &Card{Suit: Clubs, Rank: Two} // Even the weakest monster will do

	// Create a room with the monster
	room := NewRoom([]*Card{monsterCard})
	session.currentRoom = room

	// Play the monster card
	session.PlayCard(0)

	// Check loss condition
	if session.GetState() != GameStateLost {
		t.Errorf("Expected game state to be Lost after health reaches 0, but got %v", session.GetState())
	}
}
