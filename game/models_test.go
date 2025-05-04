package game

import (
	"testing"
)

func TestDeckSetup(t *testing.T) {
	// Create a new deck
	deck := NewDeck()

	// A standard deck has 52 cards
	// We should remove:
	// - Hearts: Jack, Queen, King, Ace (4 cards)
	// - Diamonds: Jack, Queen, King, Ace (4 cards)
	// Total: 8 cards removed, so deck should have 44 cards

	expectedCardCount := 52 - 8
	actualCardCount := len(deck.cards)

	if actualCardCount != expectedCardCount {
		t.Errorf("Expected deck to have %d cards after setup, but got %d", expectedCardCount, actualCardCount)
	}

	// Check that no red face cards or red aces exist in the deck
	for _, card := range deck.cards {
		if (card.Suit == Hearts || card.Suit == Diamonds) && (card.Rank >= Jack || card.Rank == Ace) {
			t.Errorf("Deck should not contain red face cards or aces, but found %s", card.String())
		}
	}

	// Check that all other cards exist in the deck
	// For non-red suits, we should have all 13 cards
	// For red suits, we should have 9 cards (2-10)

	// Count by suit
	suitCounts := make(map[Suit]int)
	for _, card := range deck.cards {
		suitCounts[card.Suit]++
	}

	// Check counts
	if suitCounts[Clubs] != 13 {
		t.Errorf("Expected 13 Clubs cards, but got %d", suitCounts[Clubs])
	}
	if suitCounts[Spades] != 13 {
		t.Errorf("Expected 13 Spades cards, but got %d", suitCounts[Spades])
	}
	if suitCounts[Hearts] != 9 {
		t.Errorf("Expected 9 Hearts cards (2-10), but got %d", suitCounts[Hearts])
	}
	if suitCounts[Diamonds] != 9 {
		t.Errorf("Expected 9 Diamonds cards (2-10), but got %d", suitCounts[Diamonds])
	}
}

func TestInitialGameState(t *testing.T) {
	session := NewGameSession()

	// Verify initial state
	if session.GetState() != GameStateInProgress {
		t.Errorf("Expected initial game state to be InProgress, got %v", session.GetState())
	}

	if session.GetPlayer().Health() != 20 {
		t.Errorf("Expected initial health to be 20, got %d", session.GetPlayer().Health())
	}

	// Verify deck setup
	// Initial deck has 44 cards (52 - 8 red face cards and aces)
	// After dealing 4 for the initial room, 40 should remain
	cardsRemaining := session.GetDeck().Remaining()
	if cardsRemaining != 40 {
		t.Errorf("Expected 40 cards remaining in deck after room setup, got %d", cardsRemaining)
	}

	// Verify room setup
	if len(session.GetCurrentRoom().Cards()) != 4 {
		t.Errorf("Expected 4 cards in initial room, got %d", len(session.GetCurrentRoom().Cards()))
	}
}
