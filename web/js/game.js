// Scoundrel Game - Core Game Logic

// API endpoints
const API = {
    BASE_URL: 'http://localhost:8080',
    NEW_GAME: '/api/games',
    GAME_STATE: '/api/games/{id}',
    PLAY_CARD: '/api/games/{id}/play/{index}',
    PLAY_CARD_WITHOUT_WEAPON: '/api/games/{id}/play-without-weapon/{index}',
    SKIP_ROOM: '/api/games/{id}/skip'
};

// Card types
const CARD_TYPE = {
    MONSTER: 0,
    WEAPON: 1,
    POTION: 2
};

// Card suits
const CARD_SUIT = {
    CLUBS: 0,    // Monsters
    DIAMONDS: 1, // Weapons
    HEARTS: 2,   // Potions
    SPADES: 3    // Monsters
};

// Game state
const GAME_STATE = {
    INITIAL: 'Initial',
    IN_PROGRESS: 'InProgress',
    WON: 'Won',
    LOST: 'Lost'
};

// Main game class
class ScoundrelGame {
    constructor() {
        this.gameId = null;
        this.gameState = null;
        this.player = {
            health: 20,
            maxHealth: 20,
            equippedWeapon: null,
            defeatedMonsters: [],
            usedPotion: false
        };
        this.room = {
            cards: [],
            completed: false
        };
        this.deck = {
            remainingCards: 0,
            previousRoomSkipped: false
        };
        this.isGameOver = false;
        this.eventListeners = {
            'gameStateChanged': [],
            'cardPlayed': [],
            'roomCompleted': [],
            'gameOver': []
        };
    }

    // Event handling
    addEventListener(event, callback) {
        if (this.eventListeners[event]) {
            this.eventListeners[event].push(callback);
        }
    }

    triggerEvent(event, data) {
        if (this.eventListeners[event]) {
            this.eventListeners[event].forEach(callback => callback(data));
        }
    }

    // API Interaction Methods
    async createNewGame() {
        try {
            const response = await fetch(`${API.BASE_URL}${API.NEW_GAME}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (!response.ok) {
                throw new Error(`Failed to create a new game: ${response.status}`);
            }
            
            const data = await response.json();
            this.gameId = data.game_id;
            
            // Now fetch initial game state
            await this.fetchGameState();
            
            // Log action
            this.addToGameLog('New game started!');
            
            return this.gameId;
        } catch (error) {
            console.error('Error creating a new game:', error);
            this.addToGameLog(`Error: ${error.message}`);
            return null;
        }
    }

    async fetchGameState() {
        if (!this.gameId) {
            throw new Error('No active game. Create a new game first.');
        }
        
        try {
            const url = `${API.BASE_URL}${API.GAME_STATE.replace('{id}', this.gameId)}`;
            const response = await fetch(url);
            
            if (!response.ok) {
                throw new Error(`Failed to fetch game state: ${response.status}`);
            }
            
            const data = await response.json();
            this.updateGameState(data);
            
            return this.gameState;
        } catch (error) {
            console.error('Error fetching game state:', error);
            this.addToGameLog(`Error: ${error.message}`);
            return null;
        }
    }

    async playCard(cardIndex, useWeapon = true) {
        if (!this.gameId) {
            throw new Error('No active game. Create a new game first.');
        }
        
        try {
            // Determine which endpoint to use based on useWeapon flag
            let endpoint = useWeapon 
                ? API.PLAY_CARD.replace('{id}', this.gameId).replace('{index}', cardIndex)
                : API.PLAY_CARD_WITHOUT_WEAPON.replace('{id}', this.gameId).replace('{index}', cardIndex);
            
            const response = await fetch(`${API.BASE_URL}${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (!response.ok) {
                throw new Error(`Failed to play card: ${response.status}`);
            }
            
            // Get the card before fetching new state
            const playedCard = this.room.cards[cardIndex];
            
            // Fetch updated game state
            await this.fetchGameState();
            
            // Generate appropriate log message
            this.logCardPlay(playedCard, useWeapon);
            
            // Trigger event
            this.triggerEvent('cardPlayed', {
                cardIndex: cardIndex,
                card: playedCard,
                useWeapon: useWeapon
            });
            
            // Check if room was completed
            if (this.room.completed) {
                this.triggerEvent('roomCompleted', {});
                this.addToGameLog('Room completed! Moving to next room...');
            }
            
            // Check for game over
            if (this.isGameOver) {
                this.triggerEvent('gameOver', {
                    won: this.gameState === GAME_STATE.WON
                });
            }
            
            return true;
        } catch (error) {
            console.error('Error playing card:', error);
            this.addToGameLog(`Error: ${error.message}`);
            return false;
        }
    }

    async skipRoom() {
        if (!this.gameId) {
            throw new Error('No active game. Create a new game first.');
        }
        
        if (this.deck.previousRoomSkipped) {
            this.addToGameLog('Cannot skip two rooms in a row!');
            return false;
        }
        
        try {
            const response = await fetch(`${API.BASE_URL}${API.SKIP_ROOM.replace('{id}', this.gameId)}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (!response.ok) {
                throw new Error(`Failed to skip room: ${response.status}`);
            }
            
            // Fetch updated game state
            await this.fetchGameState();
            
            // Log action
            this.addToGameLog('Room skipped! New room dealt.');
            
            return true;
        } catch (error) {
            console.error('Error skipping room:', error);
            this.addToGameLog(`Error: ${error.message}`);
            return false;
        }
    }

    // Game State Management
    updateGameState(data) {
        this.gameState = data.state;
        
        // Update player
        this.player.health = data.player.health;
        this.player.maxHealth = data.player.max_health;
        this.player.equippedWeapon = data.player.equipped_weapon;
        this.player.defeatedMonsters = data.player.defeated_monsters || [];
        this.player.usedPotion = data.player.used_potion;
        
        // Update room
        this.room.cards = data.room.cards || [];
        this.room.completed = data.room.completed;
        
        // Update deck
        this.deck.remainingCards = data.deck.remaining_cards;
        this.deck.previousRoomSkipped = data.deck.previous_room_skipped;
        
        // Check game over
        this.isGameOver = (data.state === GAME_STATE.WON || data.state === GAME_STATE.LOST);
        
        // Trigger state changed event
        this.triggerEvent('gameStateChanged', {
            gameState: this.gameState,
            player: this.player,
            room: this.room,
            deck: this.deck,
            isGameOver: this.isGameOver
        });
    }

    // Helper Methods
    canSkipRoom() {
        return !this.deck.previousRoomSkipped && this.room.cards.length === 4 && !this.isGameOver;
    }

    canUseWeaponAgainst(monsterCard) {
        if (!this.player.equippedWeapon) {
            return false;
        }
        
        if (this.player.defeatedMonsters.length === 0) {
            return true;
        }
        
        // Get the last defeated monster
        const lastMonster = this.player.defeatedMonsters[this.player.defeatedMonsters.length - 1];
        
        // Can only use weapon against monsters with lower value than the last one defeated
        return monsterCard.value <= lastMonster.value;
    }

    getCardTypeString(card) {
        switch (card.type) {
            case CARD_TYPE.MONSTER: return 'Monster';
            case CARD_TYPE.WEAPON: return 'Weapon';
            case CARD_TYPE.POTION: return 'Potion';
            default: return 'Card';
        }
    }

    // Logging
    logCardPlay(playedCard, useWeapon) {
        if (!playedCard) return;
        
        switch (playedCard.type) {
            case CARD_TYPE.MONSTER:
                if (useWeapon && this.player.equippedWeapon) {
                    this.addToGameLog(`Fought ${playedCard.display} (Monster) using your ${this.player.equippedWeapon.display} weapon.`);
                } else {
                    this.addToGameLog(`Fought ${playedCard.display} (Monster) barehanded! Took ${playedCard.value} damage.`);
                }
                break;
            case CARD_TYPE.WEAPON:
                this.addToGameLog(`Equipped ${playedCard.display} (Weapon).`);
                break;
            case CARD_TYPE.POTION:
                if (this.player.usedPotion) {
                    this.addToGameLog(`Used ${playedCard.display} (Potion) but it had no effect (can only use one potion per room).`);
                } else {
                    this.addToGameLog(`Used ${playedCard.display} (Potion) and healed ${playedCard.value} health.`);
                }
                break;
        }
    }

    addToGameLog(message) {
        const logContainer = document.getElementById('log-container');
        const logEntry = document.createElement('div');
        logEntry.classList.add('log-entry');
        logEntry.textContent = message;
        logContainer.appendChild(logEntry);
        
        // Scroll to bottom
        logContainer.scrollTop = logContainer.scrollHeight;
    }
}

// Create a singleton instance
const scoundrelGame = new ScoundrelGame();