package engine

import (
	"matetra/model"
)

type Game struct {
	State *model.GameState
}

// Initializes a new empty game
func New(gameID string) *Game {
	return &Game{
		State: &model.GameState{
			GameID:  gameID,
			Players: []model.Player{},
			Cards:   []model.Card{},
			Numbers: [][5]string{},
			Done:    []bool{},
			Queue:   []string{},
			Turn:    0,
		},
	}
}

// Return the index of the player whoose turn it is
func (g *Game) CurrentPlayer() int {
	n := len(g.State.Players)
	if n == 0 {
		return -1
	}
	return g.State.Turn % n
}

// Adds a new player to the game
func (g *Game) AddPlayer(name, hash string) error {
	g.State.Players = append(g.State.Players, model.Player{Name: name, Hash: hash})
	g.State.Numbers = append(g.State.Numbers, [5]string{"0", "0", "0", "0", "0"})
	g.State.Done = append(g.State.Done, false)
	return nil
}

// Adds a new card to the game
func (g *Game) AddCard(card model.Card) {
	g.State.Cards = append(g.State.Cards, card)
}

// Check if everyone has finished the turn
func (g *Game) TurnsFinished() bool {
	// Iterate over the list of players, if everyone is finished, return true, if not - false.
	return false
}

// Checks how many cards does the player have
func (g *Game) PlayerHandCount(player int) int {
	// Checks how many cards does the player have
	return 0
}

// Draws a card from deck to a player
func (g *Game) RestockCard(player int) {
	// Pick a random card from the deck
	// (If there are no more cards3 in the deck, make all the used cards into deck cards)
	// Assign that card's owner to a user
}

// Fills everyone's hands up (6 cards max)
func (g *Game) RestockCards() {
	// For every player
	// 	Count amount of cards they have.
	// 	Restock the cards to fill the hands back up. (Use the RestockCard function)
}

// Applies the card's actions (if pernament: do it, if not: only return what will happen if everything is applied)
func (g *Game) ApplyCard(cardID string, pernament bool) { // idk what to return here
	// applies a singular card
}

// Applies the card actions for all cards (if pernament: do it, if not: only return what will happen if everything is applied)
func (g *Game) ApplyCards(pernament bool) {
	// iterates the applycard for every card in the queue
}

// Ends a player's turn
func (g *Game) EndTurn() {
	// Applies All Cards (function)
	//
}

func (g *Game) RecordMove() {

}
