package engine

import (
	"matetra/cards"
	"matetra/model"
	"matetra/utils"
	"math/big"
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
			Numbers: [][5]model.Number{},
			Done:    []bool{},
			Queue:   []string{},
			Turn:    0,
		},
	}
}

func NewNumber() model.Number {
	return model.Number{
		Value: big.NewInt(0),
		Mark:  "n",
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
	nums := [5]model.Number{}
	for i := 0; i < 5; i++ {
		nums[i] = NewNumber()
	}
	g.State.Numbers = append(g.State.Numbers, nums)
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

// Makes a virtual deep copy of the game state
func (g *Game) copyState() *model.GameState {
	virtual := &model.GameState{
		GameID:  g.State.GameID,
		Players: append([]model.Player(nil), g.State.Players...),
		Cards:   append([]model.Card(nil), g.State.Cards...),
		Numbers: make([][5]model.Number, len(g.State.Numbers)),
		Done:    append([]bool(nil), g.State.Done...),
		Queue:   append([]string(nil), g.State.Queue...),
		Turn:    g.State.Turn,
	}

	for i := range g.State.Numbers {
		for j := 0; j < 5; j++ {
			orig := g.State.Numbers[i][j]
			virtual.Numbers[i][j] = model.Number{
				Mark:  orig.Mark,
				Value: new(big.Int).Set(orig.Value),
			}
		}
	}

	return virtual
}

// Applies a singular card
func (g *Game) ApplyCard(vgs *model.GameState, cardID string) {
	cards.CardFunction(vgs, cardID)

	// Remove card after applying
	for i := range vgs.Cards {
		if vgs.Cards[i].ID == cardID {
			vgs.Cards[i].Owner = -2
			vgs.Cards[i].Inputs = nil
			break
		}
	}
}

// Applies all the Cards in Queue
func (g *Game) ApplyCards(pernament bool) {
	virtual := g.copyState()

	for _, cardID := range g.State.Queue {
		g.ApplyCard(virtual, cardID)
	}

	virtual.Queue = nil

	if pernament {
		g.State = virtual
	}
}

// Ends a player's turn
func (g *Game) EndTurn() {
	// Applies All Cards (function)
	// Restocks everybody 6 cards
}

func (g *Game) RecordMove() {

}
