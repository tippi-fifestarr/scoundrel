package game

import (
	"errors"
	"math/rand"
	"time"
)

// Suit represents the card suit
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

// String returns the string representation of a suit
func (s Suit) String() string {
	return [...]string{"♣", "♦", "♥", "♠"}[s]
}

// Rank represents the card rank
type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

// String returns the string representation of a rank
func (r Rank) String() string {
	return [...]string{"", "", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}[r]
}

// CardType represents the functional type of a card
type CardType int

const (
	Monster CardType = iota
	Weapon
	Potion
)

// Card represents a playing card
type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard creates a new card with the given suit and rank
func NewCard(suit Suit, rank Rank) *Card {
	return &Card{
		Suit: suit,
		Rank: rank,
	}
}

// Value returns the numerical value of the card
func (c *Card) Value() int {
	return int(c.Rank)
}

// Type returns the functional type of the card (Monster, Weapon, or Potion)
func (c *Card) Type() CardType {
	switch c.Suit {
	case Clubs, Spades:
		return Monster
	case Diamonds:
		return Weapon
	case Hearts:
		return Potion
	default:
		return Monster
	}
}

// String returns the string representation of a card
func (c *Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

// IsRedFaceOrAce returns true if the card is a red face card or ace
func (c *Card) IsRedFaceOrAce() bool {
	return (c.Suit == Hearts || c.Suit == Diamonds) && (c.Rank >= Jack)
}

// Player represents the game player
type Player struct {
	health             int
	maxHealth          int
	equippedWeapon     *Card
	defeatedMonsters   []*Card
	usedPotionThisRoom bool
}

// NewPlayer creates a new player with the specified max health
func NewPlayer(maxHealth int) *Player {
	return &Player{
		health:             maxHealth,
		maxHealth:          maxHealth,
		defeatedMonsters:   make([]*Card, 0),
		usedPotionThisRoom: false,
	}
}

// Health returns the current health of the player
func (p *Player) Health() int {
	return p.health
}

// MaxHealth returns the maximum health of the player
func (p *Player) MaxHealth() int {
	return p.maxHealth
}

// ApplyDamage applies damage to the player
func (p *Player) ApplyDamage(amount int) error {
	p.health -= amount
	return nil
}

// Heal heals the player by the specified amount, not exceeding max health
func (p *Player) Heal(amount int) {
	p.health += amount
	// Per official rules, health cannot exceed the starting maximum (20)
	if p.health > p.maxHealth {
		p.health = p.maxHealth
	}
}

// EquipWeapon equips a weapon card
func (p *Player) EquipWeapon(card *Card) {
	p.equippedWeapon = card
	p.defeatedMonsters = make([]*Card, 0) // Clear defeated monsters
}

// EquippedWeapon returns the currently equipped weapon
func (p *Player) EquippedWeapon() *Card {
	return p.equippedWeapon
}

// AddDefeatedMonster adds a monster to the list of monsters defeated with the current weapon
func (p *Player) AddDefeatedMonster(monster *Card) {
	p.defeatedMonsters = append(p.defeatedMonsters, monster)
}

// DefeatedMonsters returns the list of monsters defeated with the current weapon
func (p *Player) DefeatedMonsters() []*Card {
	return p.defeatedMonsters
}

// CanUseWeaponAgainst checks if the current weapon can be used against a monster
func (p *Player) CanUseWeaponAgainst(monster *Card) bool {
	// If no weapon or no defeated monsters, can use weapon
	if p.equippedWeapon == nil || len(p.defeatedMonsters) == 0 {
		return true
	}

	// According to official rules, the weapon can only be used against monsters
	// with values less than or equal to the LAST monster it defeated
	// (not all previously defeated monsters as originally implemented)
	lastDefeatedMonster := p.defeatedMonsters[len(p.defeatedMonsters)-1]
	return monster.Value() <= lastDefeatedMonster.Value()
}

// UsedPotionThisRoom returns whether a potion has been used in the current room
func (p *Player) UsedPotionThisRoom() bool {
	return p.usedPotionThisRoom
}

// SetUsedPotionThisRoom sets the usedPotionThisRoom flag
func (p *Player) SetUsedPotionThisRoom(used bool) {
	p.usedPotionThisRoom = used
}

// Room represents a game room with cards
type Room struct {
	cards       []*Card
	playedCards []*Card
}

// NewRoom creates a new room with the given cards
func NewRoom(cards []*Card) *Room {
	return &Room{
		cards:       cards,
		playedCards: make([]*Card, 0, 3),
	}
}

// Cards returns the cards currently in the room
func (r *Room) Cards() []*Card {
	return r.cards
}

// PlayCard plays a card at the specified index
func (r *Room) PlayCard(index int) (*Card, error) {
	if index < 0 || index >= len(r.cards) {
		return nil, errors.New("invalid card index")
	}

	card := r.cards[index]
	r.cards = append(r.cards[:index], r.cards[index+1:]...)
	r.playedCards = append(r.playedCards, card)

	return card, nil
}

// RemainingCard returns the remaining card after playing 3 cards
func (r *Room) RemainingCard() *Card {
	if len(r.cards) == 1 {
		return r.cards[0]
	}
	return nil
}

// Completed returns true if 3 cards have been played in this room
func (r *Room) Completed() bool {
	return len(r.playedCards) == 3
}

// AllCards returns all cards in the room (played and remaining)
func (r *Room) AllCards() []*Card {
	return append(r.cards, r.playedCards...)
}

// Deck represents the collection of cards
type Deck struct {
	cards           []*Card
	prevRoomSkipped bool
}

// NewDeck creates a new deck for the game, removing red face cards and aces
func NewDeck() *Deck {
	d := &Deck{
		cards:           make([]*Card, 0, 44), // 52 - 8 (red face cards and aces)
		prevRoomSkipped: false,
	}

	// Initialize deck with all cards except red face cards and aces
	for s := Clubs; s <= Spades; s++ {
		for r := Two; r <= Ace; r++ {
			card := NewCard(Suit(s), Rank(r))
			if !card.IsRedFaceOrAce() {
				d.cards = append(d.cards, card)
			}
		}
	}

	return d
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw draws the specified number of cards from the deck
func (d *Deck) Draw(count int) ([]*Card, error) {
	if count > len(d.cards) {
		return nil, errors.New("not enough cards in deck")
	}

	drawn := d.cards[:count]
	d.cards = d.cards[count:]
	return drawn, nil
}

// AddToBottom adds cards to the bottom of the deck
func (d *Deck) AddToBottom(cards []*Card) {
	d.cards = append(d.cards, cards...)
}

// Remaining returns the number of cards remaining in the deck
func (d *Deck) Remaining() int {
	return len(d.cards)
}

// PrevRoomSkipped returns whether the previous room was skipped
func (d *Deck) PrevRoomSkipped() bool {
	return d.prevRoomSkipped
}

// SetPrevRoomSkipped sets the prevRoomSkipped flag
func (d *Deck) SetPrevRoomSkipped(skipped bool) {
	d.prevRoomSkipped = skipped
}
