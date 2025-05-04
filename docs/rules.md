# Scoundrel Game Rules and Test Cases

This document outlines the rules of Scoundrel with corresponding pseudocode test cases to ensure correct implementation. It also includes questions about potentially ambiguous rules that need clarification.

## Official Rules Overview

Scoundrel is a single player rogue-like card game by Zach Gage and Kurt Bieg. These test cases and rules are based on version 1.0 (August 15th, 2011).

## Setup Rules

### Deck Construction
- **Rule**: A standard 52-card deck is used, with red face cards and red aces removed
- **Details**: Remove Jack, Queen, King, and Ace of Hearts and Diamonds (8 cards total)
- **Test Case**:
  ```
  deck = new Deck()
  assert deck.totalCards == 44
  assert deck.containsNoRedFaceCards()
  assert deck.containsNoRedAces()
  ```

### Initial Health
- **Rule**: Player starts with 20 health points
- **Test Case**:
  ```
  game = new Game()
  assert game.player.health == 20
  assert game.player.maxHealth == 20
  ```

### Room Creation
- **Rule**: A room is created by dealing 4 cards face up from the deck
- **Test Case**:
  ```
  game = new Game()
  assert game.currentRoom.cards.length == 4
  assert game.deck.remainingCards == 40  // 44 - 4
  ```

## Card Types and Values

### Card Values
- **Rule**: Cards have a value corresponding to their rank, with Ace high
- **Details**:
  - Number cards (2-10): Value equals their number
  - Jack = 11, Queen = 12, King = 13, Ace = 14
- **Test Case**:
  ```
  assert new Card(Spades, Two).value == 2
  assert new Card(Clubs, Ten).value == 10
  assert new Card(Spades, Jack).value == 11
  assert new Card(Clubs, Queen).value == 12
  assert new Card(Spades, King).value == 13
  assert new Card(Clubs, Ace).value == 14
  ```

### Monster Cards
- **Rule**: The 26 Clubs and Spades in the deck are Monsters with damage equal to their value
- **Test Case**:
  ```
  assert new Card(Clubs, Ten).type == Monster
  assert new Card(Spades, Jack).type == Monster
  assert new Card(Clubs, Ten).damage == 10
  assert new Card(Spades, Jack).damage == 11
  ```

### Weapon Cards
- **Rule**: The 9 Diamonds in the deck are Weapons with damage reduction equal to their value
- **Rule**: All weapons are binding - picking one up means you must equip it and discard your previous weapon
- **Test Case**:
  ```
  assert new Card(Diamonds, Ten).type == Weapon
  assert new Card(Diamonds, Ten).value == 10
  
  // Test binding behavior
  player.equipWeapon(new Card(Diamonds, Five))
  oldWeapon = player.equippedWeapon
  player.equipWeapon(new Card(Diamonds, Eight))
  assert player.equippedWeapon != oldWeapon
  ```

### Potion Cards
- **Rule**: The 9 Hearts in the deck are Health Potions with healing equal to their value
- **Rule**: Health may not exceed the starting 20 health
- **Test Case**:
  ```
  assert new Card(Hearts, Ten).type == Potion
  assert new Card(Hearts, Ten).healingValue == 10
  
  // Test max health cap
  player.health = 15
  player.heal(new Card(Hearts, Ten))
  assert player.health == 20  // Not 25, capped at max
  ```

## Room Mechanics

### Playing Cards
- **Rule**: Player must play 3 cards from each room in chosen order
- **Test Case**:
  ```
  game = new Game()
  initialRoomCards = game.currentRoom.cards.length  // Should be 4
  
  game.playCard(0)
  assert game.currentRoom.cards.length == initialRoomCards - 1
  
  game.playCard(0)
  assert game.currentRoom.cards.length == initialRoomCards - 2
  
  game.playCard(0)
  assert game.currentRoom.cards.length == 0
  assert game.currentRoom.completed == true
  ```

### Room Progression
- **Rule**: The remaining card starts the next room with 3 new cards
- **Test Case**:
  ```
  game = new Game()
  
  // Play 3 cards to complete the room
  remainingCardBeforeCompletion = game.currentRoom.cards[3]
  game.playCard(0)
  game.playCard(0)
  game.playCard(0)
  
  // Check that remaining card is used in the next room
  assert game.currentRoom.cards[0] == remainingCardBeforeCompletion
  assert game.currentRoom.cards.length == 4
  ```

### Room Avoidance (Skipping)
- **Rule**: A room may be avoided, but never two rooms in a row
- **Rule**: To avoid a room, all four cards are placed at the bottom of the deck
- **Test Case**:
  ```
  game = new Game()
  
  // Avoid first room
  initialCardsBefore = game.currentRoom.cards
  game.skipRoom()
  assert game.currentRoom.cards != initialCardsBefore
  
  // Try avoiding again (should fail)
  try {
    game.skipRoom()
    assert false  // Should not reach here
  } catch (error) {
    assert error.message.contains("cannot avoid two rooms in a row")
  }
  ```

## Combat Mechanics

### Basic Monster Damage
- **Rule**: When fighting a monster barehanded, player takes full damage equal to monster's value
- **Test Case**:
  ```
  game = new Game()
  monsterCard = new Card(Clubs, Ten)  // Monster with value 10
  
  initialHealth = game.player.health
  game.handleMonsterWithoutWeapon(monsterCard)
  assert game.player.health == initialHealth - 10
  ```

### Weapon Usage
- **Rule**: When fighting with a weapon, damage is calculated as (monster value - weapon value)
- **Rule**: If weapon value exceeds monster value, player takes no damage
- **Test Case**:
  ```
  game = new Game()
  weaponCard = new Card(Diamonds, Eight)  // Weapon with value 8
  monsterCard = new Card(Clubs, Ten)      // Monster with value 10
  
  game.player.equipWeapon(weaponCard)
  initialHealth = game.player.health
  game.handleMonsterWithWeapon(monsterCard)
  assert game.player.health == initialHealth - 2  // 10 - 8 = 2 damage
  
  // Test no damage when weapon > monster
  strongWeapon = new Card(Diamonds, Ten)
  weakMonster = new Card(Clubs, Three)
  
  game.player.equipWeapon(strongWeapon)
  initialHealth = game.player.health
  game.handleMonsterWithWeapon(weakMonster)
  assert game.player.health == initialHealth  // No damage taken
  ```

### Optional Weapon Use
- **Rule**: Player may choose to fight without a weapon even if one is equipped
- **Test Case**:
  ```
  game = new Game()
  weaponCard = new Card(Diamonds, Eight)  // Weapon with value 8
  monsterCard = new Card(Clubs, Ten)      // Monster with value 10
  
  game.player.equipWeapon(weaponCard)
  initialHealth = game.player.health
  
  // Choose not to use weapon
  game.handleMonsterWithoutWeapon(monsterCard)
  assert game.player.health == initialHealth - 10  // Full damage
  ```

### Weapon Restriction Rule
- **Rule**: Once a weapon has been used on a monster, it can only be used to slay monsters of a value less than or equal to the previous monster it slayed
- **Test Case**:
  ```
  game = new Game()
  weaponCard = new Card(Diamonds, Five)       // Weapon with value 5
  firstMonster = new Card(Clubs, Ten)         // Monster with value 10
  weakerMonster = new Card(Spades, Six)       // Monster with value 6
  evenWeakerMonster = new Card(Clubs, Five)   // Monster with value 5
  strongerMonster = new Card(Spades, Seven)   // Monster with value 7
  
  // Equip weapon
  game.player.equipWeapon(weaponCard)
  
  // Defeat first monster
  game.handleMonsterWithWeapon(firstMonster)
  
  // Fight weaker monster (should work)
  initialHealth = game.player.health
  game.handleMonsterWithWeapon(weakerMonster)
  assert game.player.health == initialHealth - (weakerMonster.value - weaponCard.value)
  
  // Fight even weaker monster (should work)
  initialHealth = game.player.health
  game.handleMonsterWithWeapon(evenWeakerMonster)
  assert game.player.health == initialHealth - (evenWeakerMonster.value - weaponCard.value)
  
  // Fight stronger monster than the most recent (should NOT work with weapon)
  initialHealth = game.player.health
  // This should force barehanded combat automatically
  game.handleMonsterWithWeapon(strongerMonster)
  assert game.player.health == initialHealth - strongerMonster.value  // Full damage, weapon unusable
  ```

### Weapon Replacement
- **Rule**: Equipping a new weapon discards the previous weapon along with all monsters slain by it
- **Test Case**:
  ```
  game = new Game()
  weaponCard1 = new Card(Diamonds, Five)
  weaponCard2 = new Card(Diamonds, Eight)
  monsterCard = new Card(Clubs, Ten)
  
  // Equip first weapon and defeat a monster
  game.player.equipWeapon(weaponCard1)
  game.handleMonsterWithWeapon(monsterCard)
  assert game.player.defeatedMonsters.length == 1
  
  // Equip new weapon, history should reset
  game.player.equipWeapon(weaponCard2)
  assert game.player.defeatedMonsters.length == 0
  ```

## Potion Mechanics

### Health Potions
- **Rule**: Health potions restore health equal to their value, up to the maximum of 20
- **Test Case**:
  ```
  game = new Game()
  game.player.health = 10  // Reduce health
  potionCard = new Card(Hearts, Seven)
  
  game.handlePotion(potionCard)
  assert game.player.health == 17  // 10 + 7
  ```

### Multiple Potions
- **Rule**: Only the first potion in a room has an effect, others are simply discarded
- **Test Case**:
  ```
  game = new Game()
  game.player.health = 10  // Reduce health
  potionCard1 = new Card(Hearts, Five)
  potionCard2 = new Card(Hearts, Seven)
  
  game.handlePotion(potionCard1)
  assert game.player.health == 15  // 10 + 5
  
  game.handlePotion(potionCard2)
  assert game.player.health == 15  // No change
  ```

### Maximum Health Cap
- **Rule**: Health may not exceed the starting value of 20
- **Test Case**:
  ```
  game = new Game()
  game.player.health = 18  // Close to max
  potionCard = new Card(Hearts, Seven)
  
  game.handlePotion(potionCard)
  assert game.player.health == 20  // Capped at max, not 25
  ```

## Win/Loss Conditions

### Winning
- **Rule**: Player wins by exhausting the dungeon deck
- **Test Case**:
  ```
  game = new Game()
  
  // Play through all cards
  while (game.deck.remaining > 0) {
    game.playCard(0)  // Simplification: just play first card each time
  }
  
  assert game.isWon() == true
  ```

### Losing
- **Rule**: Player loses if health reaches 0 or below
- **Test Case**:
  ```
  game = new Game()
  
  // Take enough damage to lose
  game.player.health = 5
  monsterCard = new Card(Clubs, Ten)
  
  game.handleMonsterWithoutWeapon(monsterCard)
  assert game.player.health <= 0
  assert game.isLost() == true
  ```

## Scoring System

### Losing Score
- **Rule**: If you lose, your score is negative - the sum of all remaining monsters in the dungeon
- **Test Case**:
  ```
  game = new Game()
  
  // Lose the game
  game.player.health = 0
  
  // Add monsters to remaining deck
  game.deck.cards = [
    new Card(Clubs, Five),
    new Card(Spades, Ten),
    new Card(Hearts, Seven), // Potion, not counted
    new Card(Clubs, King)
  ]
  
  assert game.calculateScore() == -(5 + 10 + 13)  // Negative sum of monster values
  ```

### Winning Score
- **Rule**: If you win, your score is your remaining health
- **Rule**: If your health is 20 and your last card was a health potion, your score is 20 + potion value
- **Test Case**:
  ```
  game = new Game()
  
  // Win with some health remaining
  game.player.health = 15
  game.deck.cards = []  // Empty deck means win
  
  assert game.calculateScore() == 15
  
  // Win with full health and last card was potion
  game.player.health = 20
  game.lastCardPlayed = new Card(Hearts, Seven)
  
  assert game.calculateScore() == 27  // 20 + 7
  ```

## Open Questions and Implementation Notes

1. **Optional Weapon Usage**:
   - The rules clearly state that players may choose to fight with or without a weapon
   - Our current implementation doesn't allow for this choice explicitly
   - Need to modify the engine to support this feature

2. **Weapon Restriction Rule**:
   - The official rules indicate that the restriction is based on the most recently defeated monster
   - Our implementation currently checks against ALL previously defeated monsters
   - Need to update the CanUseWeaponAgainst method to only check the most recent monster

3. **Scoring System**:
   - Not currently implemented in our game engine
   - Need to add a calculateScore method based on the official rules

4. **Health Cap**:
   - Official rules confirm that health cannot exceed 20
   - Our implementation correctly enforces this already

5. **Win Event Emission**:
   - Need to determine how to properly notify the player of victory/defeat
   - Consider adding special game state transitions and UI notifications