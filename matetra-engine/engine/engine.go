package engine

import (
	"matetra/cards"
	"matetra/model"
	"matetra/utils"
	"math/rand"
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

// Loads the card deck into game state
func (g *Game) LoadCards() {
	deck := utils.Must(cards.LoadCardsFromCSV(utils.DECK_PATH))
	g.State.Cards = append(g.State.Cards, deck...)
}

// Check if everyone has finished the turn
func (g *Game) TurnsFinished() bool {
	for _, done := range g.State.Done {
		if !done {
			return false
		}
	}
	return true
}

// Checks how many cards does the player have
func (g *Game) PlayerHandCount(player int) int {
	count := 0
	for _, card := range g.State.Cards {
		if card.Owner == player {
			count++
		}
	}
	return count
}

// Fills everyone's hands up (6 cards max)
func (g *Game) RestockCards() {
	for p := range g.State.Players {
		hand := g.PlayerHandCount(p)

		for hand < 6 {
			// Build a deck
			deck := []int{}
			for i, c := range g.State.Cards {
				if c.Owner == -1 {
					deck = append(deck, i)
				}
			}

			// If deck is empty, recycle used cards and build a deck
			if len(deck) == 0 {
				for i := range g.State.Cards {
					if g.State.Cards[i].Owner == -2 {
						g.State.Cards[i].Owner = -1
					}
				}

				for i, c := range g.State.Cards {
					if c.Owner == -1 {
						deck = append(deck, i)
					}
				}
			}

			// Choose a card from a deck
			idx := deck[rand.Intn(len(deck))]
			g.State.Cards[idx].Owner = p
			hand++
		}
	}
}

// Applies the card's actions (if pernament: do it, if not: only return what will happen if everything is applied)
func (g *Game) ApplyCard(cardID string, pernament bool) { // idk what to return here
	// Make a new VirtualGameState which copies everything from a game state
	virtualGameState := &model.GameState{
		GameID:  g.State.GameID,
		Players: append([]model.Player(nil), g.State.Players...),
		Cards:   append([]model.Card(nil), g.State.Cards...),
		Numbers: make([][5]string, len(g.State.Numbers)),
		Done:    append([]bool(nil), g.State.Done...),
		Queue:   append([]string(nil), g.State.Queue...),
		Turn:    g.State.Turn,
	}
	copy(virtualGameState.Numbers, g.State.Numbers)
	// cards.CardFunction(virtualGameState.Numbers, g.State.Numbers)
	// Call the CardFunction
	// If pernament, replace the GameState with the VirtualGameState
	// Clean up the card data and set it as used
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
