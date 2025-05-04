package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tippi-fifestarr/scoundrel/game"
)

func main() {
	fmt.Println("Scoundrel Card Game CLI")
	fmt.Println("=======================")
	fmt.Println("Starting new game...")

	// Create a new game session
	session := game.NewGameSession()
	reader := bufio.NewReader(os.Stdin)

	// Game loop
	for !session.IsGameOver() {
		// Display game state
		displayGameState(session)

		// Get player action
		action, err := getPlayerAction(reader, session)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		// Process action
		executeAction(action, reader, session)
	}

	// Game over
	displayGameState(session)
	if session.GetState() == game.GameStateWon {
		fmt.Println("Congratulations! You won!")
	} else {
		fmt.Println("Game over! You lost.")
	}
}

func displayGameState(session *game.GameSession) {
	fmt.Println("\n--------------------------------------------------")
	fmt.Printf("Health: %d/%d\n", session.GetPlayer().Health(), session.GetPlayer().MaxHealth())

	// Display equipped weapon
	if weapon := session.GetPlayer().EquippedWeapon(); weapon != nil {
		fmt.Printf("Equipped Weapon: %s (Value: %d)\n", weapon.String(), weapon.Value())

		if len(session.GetPlayer().DefeatedMonsters()) > 0 {
			fmt.Print("Defeated monsters: ")
			for i, monster := range session.GetPlayer().DefeatedMonsters() {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(monster.String())
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No weapon equipped")
	}

	// Add potion usage status for clarity
	if session.GetPlayer().UsedPotionThisRoom() {
		fmt.Println("Potion already used in this room (only one effective potion per room)")
	}

	// Display current room
	if room := session.GetCurrentRoom(); room != nil {
		fmt.Println("\nCurrent Room:")
		for i, card := range room.Cards() {
			fmt.Printf("[%d] %s ", i, card.String())

			switch card.Type() {
			case game.Monster:
				fmt.Printf("(Monster, Damage: %d)", card.Value())
			case game.Weapon:
				fmt.Printf("(Weapon, Value: %d)", card.Value())
			case game.Potion:
				fmt.Printf("(Potion, Heal: %d)", card.Value())
			}
			fmt.Println()
		}
	}

	// Display deck info
	fmt.Printf("\nCards remaining in dungeon: %d\n", session.GetDeck().Remaining())
	if session.GetDeck().PrevRoomSkipped() {
		fmt.Println("You skipped the previous room, you cannot skip this one.")
	}
	fmt.Println("--------------------------------------------------")
}

func getPlayerAction(reader *bufio.Reader, session *game.GameSession) (string, error) {
	// Display options
	fmt.Println("\nActions:")
	for i, card := range session.GetCurrentRoom().Cards() {
		fmt.Printf("[%d] Play card %s\n", i, card.String())
	}

	// Skip room option - only show if the previous room wasn't skipped AND no cards have been played yet
	currentRoom := session.GetCurrentRoom()
	if !session.GetDeck().PrevRoomSkipped() && len(currentRoom.Cards()) == 4 { // Original room has 4 cards, so no cards played yet
		fmt.Println("[s] Skip this room")
	}

	fmt.Println("[q] Quit game")
	fmt.Print("\nEnter your choice: ")

	// Read input
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func executeAction(action string, reader *bufio.Reader, session *game.GameSession) {
	// Check for quit
	if action == "q" {
		fmt.Println("Quitting game...")
		os.Exit(0)
	}

	// Check for skip room
	if action == "s" {
		// Only allow skipping if the previous room wasn't skipped AND no cards have been played yet
		currentRoom := session.GetCurrentRoom()
		if session.GetDeck().PrevRoomSkipped() {
			fmt.Println("You cannot skip two rooms in a row!")
			return
		}

		if len(currentRoom.Cards()) < 4 { // If cards have been played, the room has less than 4 cards
			fmt.Println("You cannot skip a room after playing cards!")
			return
		}

		err := session.SkipRoom()
		if err != nil {
			fmt.Printf("Error skipping room: %s\n", err)
		} else {
			fmt.Println("Room skipped! New room dealt.")
		}
		return
	}

	// Try to play a card
	index, err := strconv.Atoi(action)
	if err != nil {
		fmt.Println("Invalid action! Please try again.")
		return
	}

	// Check if the index is valid
	if index < 0 || index >= len(session.GetCurrentRoom().Cards()) {
		fmt.Println("Invalid card index! Please try again.")
		return
	}

	// Get the card before playing it to provide better feedback and check if it's a monster
	card := session.GetCurrentRoom().Cards()[index]

	// If it's a monster and player has a weapon, ask if they want to use it
	if card.Type() == game.Monster && session.GetPlayer().EquippedWeapon() != nil {
		weapon := session.GetPlayer().EquippedWeapon()

		// Check if weapon can be used against this monster
		canUseWeapon := session.GetPlayer().CanUseWeaponAgainst(card)
		if canUseWeapon {
			damage := card.Value() - weapon.Value()
			if damage < 0 {
				damage = 0
			}

			fmt.Printf("\nYou're facing a monster with value %d. You have a weapon with value %d.\n", card.Value(), weapon.Value())
			fmt.Printf("Using your weapon would result in %d damage.\n", damage)
			fmt.Printf("Fighting barehanded would result in %d damage.\n", card.Value())
			fmt.Print("Do you want to use your weapon? (y/n): ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if strings.ToLower(input) == "y" || strings.ToLower(input) == "yes" {
				// Use weapon
				err = session.PlayCard(index)
			} else {
				// Fight barehanded
				err = session.PlayCardWithoutWeapon(index)
			}
		} else {
			fmt.Printf("\nYour weapon (%s) can't be used against this monster because it's stronger than the last monster you defeated.\n", weapon.String())
			fmt.Printf("You'll take full damage of %d from this monster.\n", card.Value())

			// Automatically fight barehanded
			err = session.PlayCardWithoutWeapon(index)
		}
	} else {
		// For non-monster cards or when player has no weapon, just play normally
		err = session.PlayCard(index)
	}

	if err != nil {
		fmt.Printf("Error playing card: %s\n", err)
		return
	}

	// Provide feedback based on card type
	switch card.Type() {
	case game.Monster:
		fmt.Printf("Fought a monster with value %d.\n", card.Value())
	case game.Weapon:
		fmt.Printf("Equipped a weapon with value %d.\n", card.Value())
	case game.Potion:
		if !session.GetPlayer().UsedPotionThisRoom() {
			fmt.Printf("Used a potion with value %d and restored health.\n", card.Value())
		} else {
			fmt.Printf("Used a potion with value %d but it had no effect (only one effective potion per room).\n", card.Value())
		}
	}
}
