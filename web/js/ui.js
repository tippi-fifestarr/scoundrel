// Scoundrel Game - UI Handling

// Main UI controller
class ScoundrelUI {
    constructor(game) {
        this.game = game;
        this.setupEventListeners();
        this.setupGameListeners();
    }

    // Set up DOM event listeners
    setupEventListeners() {
        // Button listeners
        document.getElementById('new-game-btn').addEventListener('click', () => this.startNewGame());
        document.getElementById('skip-room-btn').addEventListener('click', () => this.skipRoom());
        document.getElementById('use-weapon-btn').addEventListener('click', () => this.resolveCombat(true));
        document.getElementById('fight-barehanded-btn').addEventListener('click', () => this.resolveCombat(false));
        document.getElementById('restart-game-btn').addEventListener('click', () => this.startNewGame());

        // Initialize skip room button state
        this.updateSkipRoomButton();
    }

    // Set up game event listeners
    setupGameListeners() {
        // Listen for game state changes
        this.game.addEventListener('gameStateChanged', (data) => {
            this.renderGameState(data);
            this.updateSkipRoomButton();
        });

        // Listen for card played events
        this.game.addEventListener('cardPlayed', (data) => {
            this.animateCardPlay(data.cardIndex, data.card);
        });

        // Listen for room completed events
        this.game.addEventListener('roomCompleted', () => {
            setTimeout(() => {
                this.renderRoom();
            }, 500);
        });

        // Listen for game over events
        this.game.addEventListener('gameOver', (data) => {
            this.showGameOverScreen(data.won);
        });
    }

    // Start a new game
    async startNewGame() {
        // Clear the UI
        this.clearGameUI();
        
        // Show loading state
        this.setLoadingState(true);
        
        // Create a new game via API
        await this.game.createNewGame();
        
        // Render initial state
        this.renderGameState({
            gameState: this.game.gameState,
            player: this.game.player,
            room: this.game.room,
            deck: this.game.deck
        });
        
        // Hide loading state
        this.setLoadingState(false);
        
        // Hide game over modal if visible
        document.getElementById('gameover-modal').style.display = 'none';
    }

    // Skip the current room
    async skipRoom() {
        if (!this.game.canSkipRoom()) {
            return;
        }
        
        // Disable button during processing
        const skipButton = document.getElementById('skip-room-btn');
        skipButton.disabled = true;
        
        // Call API
        await this.game.skipRoom();
        
        // Re-enable button (will be hidden if no longer applicable)
        skipButton.disabled = false;
    }

    // Handle card click
    async handleCardClick(cardIndex) {
        const card = this.game.room.cards[cardIndex];
        if (!card) return;
        
        // If this is a monster and player has a weapon, show combat options
        if (card.type === CARD_TYPE.MONSTER && 
            this.game.player.equippedWeapon && 
            this.game.canUseWeaponAgainst(card)) {
                
            this.showCombatModal(cardIndex, card);
        } else {
            // Otherwise just play the card normally
            await this.game.playCard(cardIndex);
        }
    }

    // Show combat choice modal
    showCombatModal(cardIndex, monsterCard) {
        this.currentCombat = { cardIndex, monsterCard };
        
        const modal = document.getElementById('combat-modal');
        const description = document.getElementById('combat-description');
        
        const weaponValue = this.game.player.equippedWeapon.value;
        const monsterValue = monsterCard.value;
        const reducedDamage = Math.max(0, monsterValue - weaponValue);
        
        description.innerHTML = `
            <p>You're facing a ${monsterCard.display} (Monster, Value: ${monsterValue}).</p>
            <p>Your weapon is ${this.game.player.equippedWeapon.display} (Value: ${weaponValue}).</p>
            <p>Using your weapon would result in <strong>${reducedDamage}</strong> damage.</p>
            <p>Fighting barehanded would result in <strong>${monsterValue}</strong> damage.</p>
            <p>How do you want to fight?</p>
        `;
        
        modal.style.display = 'flex';
    }

    // Resolve combat after player choice
    async resolveCombat(useWeapon) {
        // Hide modal
        document.getElementById('combat-modal').style.display = 'none';
        
        // Play the card with the chosen combat approach
        if (this.currentCombat) {
            await this.game.playCard(this.currentCombat.cardIndex, useWeapon);
            this.currentCombat = null;
        }
    }

    // Show game over screen
    showGameOverScreen(won) {
        const modal = document.getElementById('gameover-modal');
        const title = document.getElementById('gameover-title');
        const message = document.getElementById('gameover-message');
        
        if (won) {
            title.textContent = 'Victory!';
            message.textContent = `Congratulations! You defeated the dungeon with ${this.game.player.health} health remaining.`;
        } else {
            title.textContent = 'Game Over';
            message.textContent = 'You were defeated by the dungeon. Better luck next time!';
        }
        
        modal.style.display = 'flex';
    }

    // Render the full game state
    renderGameState(data) {
        this.renderHealth();
        this.renderWeapon();
        this.renderDefeatedMonsters();
        this.renderDeck();
        this.renderRoom();
    }

    // Render health display
    renderHealth() {
        const health = this.game.player.health;
        const maxHealth = this.game.player.maxHealth;
        const healthPercent = Math.max(0, Math.min(100, (health / maxHealth) * 100));
        
        // Update health bar
        const healthFill = document.getElementById('health-fill');
        healthFill.style.width = `${healthPercent}%`;
        
        // Update health text
        const healthValue = document.getElementById('health-value');
        healthValue.textContent = `${health}/${maxHealth}`;
        
        // Color coding
        if (healthPercent > 60) {
            healthFill.style.backgroundColor = '#2ecc71'; // Green
        } else if (healthPercent > 30) {
            healthFill.style.backgroundColor = '#f39c12'; // Orange
        } else {
            healthFill.style.backgroundColor = '#e74c3c'; // Red
        }
    }

    // Render weapon display
    renderWeapon() {
        const weaponSlot = document.getElementById('weapon-slot');
        weaponSlot.innerHTML = '';
        
        if (this.game.player.equippedWeapon) {
            const weapon = this.game.player.equippedWeapon;
            weaponSlot.appendChild(this.createCardElement(weapon));
        }
    }

    // Render defeated monsters
    renderDefeatedMonsters() {
        const monstersContainer = document.getElementById('monsters-container');
        monstersContainer.innerHTML = '';
        
        if (this.game.player.defeatedMonsters && this.game.player.defeatedMonsters.length > 0) {
            this.game.player.defeatedMonsters.forEach(monster => {
                monstersContainer.appendChild(this.createCardElement(monster));
            });
        } else {
            const emptyMessage = document.createElement('p');
            emptyMessage.textContent = 'No monsters defeated yet';
            emptyMessage.classList.add('empty-message');
            monstersContainer.appendChild(emptyMessage);
        }
    }

    // Render deck information
    renderDeck() {
        const deckCount = document.getElementById('deck-count');
        deckCount.textContent = this.game.deck.remainingCards;
    }

    // Render room cards
    renderRoom() {
        // Clear all card slots
        for (let i = 0; i < 4; i++) {
            const slot = document.getElementById(`card-slot-${i}`);
            slot.innerHTML = '';
        }
        
        // Add cards to appropriate slots
        this.game.room.cards.forEach((card, index) => {
            const slot = document.getElementById(`card-slot-${index}`);
            const cardElement = this.createCardElement(card);
            
            // Add click event
            cardElement.addEventListener('click', () => this.handleCardClick(index));
            
            slot.appendChild(cardElement);
        });
    }

    // Create a card DOM element
    createCardElement(card) {
        const cardElement = document.createElement('div');
        cardElement.classList.add('card');
        
        // Add appropriate class based on card type
        switch (card.type) {
            case CARD_TYPE.MONSTER:
                cardElement.classList.add('card-monster');
                break;
            case CARD_TYPE.WEAPON:
                cardElement.classList.add('card-weapon');
                break;
            case CARD_TYPE.POTION:
                cardElement.classList.add('card-potion');
                break;
        }
        
        // Create card structure
        const cardTop = document.createElement('div');
        cardTop.classList.add('card-top');
        
        const cardValue = document.createElement('div');
        cardValue.classList.add('card-value');
        cardValue.textContent = card.display.replace(/[♣♦♥♠]/g, '');
        
        const cardSuit = document.createElement('div');
        cardSuit.classList.add('card-suit');
        
        // Add suit text
        const suitChar = card.display.match(/[♣♦♥♠]/)[0];
        cardSuit.textContent = suitChar;
        
        // Add suit-specific class for coloring
        switch (suitChar) {
            case '♣':
                cardSuit.classList.add('suit-clubs');
                break;
            case '♦':
                cardSuit.classList.add('suit-diamonds');
                break;
            case '♥':
                cardSuit.classList.add('suit-hearts');
                break;
            case '♠':
                cardSuit.classList.add('suit-spades');
                break;
        }
        
        cardTop.appendChild(cardValue);
        cardTop.appendChild(cardSuit);
        
        const cardCenter = document.createElement('div');
        cardCenter.classList.add('card-center');
        
        switch (card.type) {
            case CARD_TYPE.MONSTER:
                cardCenter.textContent = '👾';
                break;
            case CARD_TYPE.WEAPON:
                cardCenter.textContent = '⚔️';
                break;
            case CARD_TYPE.POTION:
                cardCenter.textContent = '⚗️';
                break;
        }
        
        const cardBottom = document.createElement('div');
        cardBottom.classList.add('card-bottom');
        
        const bottomValue = document.createElement('div');
        bottomValue.classList.add('card-value');
        bottomValue.textContent = card.display.replace(/[♣♦♥♠]/g, '');
        
        const bottomSuit = document.createElement('div');
        bottomSuit.classList.add('card-suit');
        bottomSuit.textContent = suitChar;
        
        // Add the same suit class to bottom suit
        switch (suitChar) {
            case '♣':
                bottomSuit.classList.add('suit-clubs');
                break;
            case '♦':
                bottomSuit.classList.add('suit-diamonds');
                break;
            case '♥':
                bottomSuit.classList.add('suit-hearts');
                break;
            case '♠':
                bottomSuit.classList.add('suit-spades');
                break;
        }
        
        cardBottom.appendChild(bottomValue);
        cardBottom.appendChild(bottomSuit);
        
        cardElement.appendChild(cardTop);
        cardElement.appendChild(cardCenter);
        cardElement.appendChild(cardBottom);
        
        return cardElement;
    }

    // Animate a card being played
    animateCardPlay(cardIndex, card) {
        // Implementation would depend on desired animation
        // For now, we'll just hide the card and then update the UI
        const slot = document.getElementById(`card-slot-${cardIndex}`);
        const cardElement = slot.querySelector('.card');
        
        if (cardElement) {
            cardElement.style.transition = 'transform 0.5s, opacity 0.5s';
            cardElement.style.opacity = '0';
            cardElement.style.transform = 'translateY(-50px)';
            
            setTimeout(() => {
                this.renderRoom();
            }, 500);
        }
    }

    // Set loading state
    setLoadingState(isLoading) {
        const buttons = document.querySelectorAll('button');
        buttons.forEach(button => {
            button.disabled = isLoading;
        });
        
        // Could add spinner or other loading indicators here
    }

    // Clear the game UI
    clearGameUI() {
        // Clear room
        for (let i = 0; i < 4; i++) {
            const slot = document.getElementById(`card-slot-${i}`);
            slot.innerHTML = '';
        }
        
        // Clear weapon
        document.getElementById('weapon-slot').innerHTML = '';
        
        // Clear defeated monsters
        document.getElementById('monsters-container').innerHTML = '';
        
        // Reset health
        document.getElementById('health-fill').style.width = '100%';
        document.getElementById('health-value').textContent = '20/20';
        
        // Clear log
        document.getElementById('log-container').innerHTML = '';
    }

    // Update skip room button state
    updateSkipRoomButton() {
        const skipButton = document.getElementById('skip-room-btn');
        skipButton.disabled = !this.game.canSkipRoom();
    }
}

// Initialize the UI when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const ui = new ScoundrelUI(scoundrelGame);
    
    // Auto-start a new game when the page loads
    ui.startNewGame();
});