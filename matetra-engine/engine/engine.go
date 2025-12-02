package engine

import (
	"fmt"
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
			Numbers: make([][5]model.Number, 0),
			Done:    make([]bool, 0),
			Queue:   make([]string, 0),
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

func NewNumberRow() (row [5]model.Number) {
	for i := range row {
		row[i] = NewNumber()
	}
	return
}

// Adds a new player to the game
func (g *Game) AddPlayer(name, hash string) error {
	// Adds a new player object
	g.State.Players = append(g.State.Players, model.Player{
		Name: name,
		Hash: hash,
	})

	// Adds a new number row
	g.State.Numbers = append(g.State.Numbers, NewNumberRow())

	// Adds a new 'not done' flag
	g.State.Done = append(g.State.Done, false)
	return nil
}

// Return the index of the player whoose turn it is
func (g *Game) CurrentPlayer() int {
	n := len(g.State.Players)
	if n == 0 {
		return -1
	}
	return g.State.Turn % n
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

// Fills everyone's hands up (6 cards max) (needs optimisation)
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
func (g *Game) CopyState() *model.GameState {
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
func (g *Game) ApplyCard(vgs *model.GameState, cardID string) error {
	err := cards.CardFunction(vgs, cardID)
	if err != nil {
		return err
	}

	// Remove card after applying
	for i := range vgs.Cards {
		if vgs.Cards[i].ID == cardID {
			vgs.Cards[i].Owner = -2
			vgs.Cards[i].Inputs = nil
			break
		}
	}

	return nil
}

// Applies all the Cards in Queue
func (g *Game) ApplyCards(pernament bool) (*model.GameState, error) {
	virtual := g.CopyState()

	for _, cardID := range g.State.Queue {
		err := g.ApplyCard(virtual, cardID)
		if err != nil {
			return nil, err
		}
	}

	virtual.Queue = nil

	if pernament {
		g.State = virtual
	}

	return virtual, nil
}

func (g *Game) IncrementPlayer() {
	g.State.Turn += 1
}

// Ends a player's turn
func (g *Game) EndTurn() {
	g.ApplyCards(true)
	g.RestockCards()
	g.IncrementPlayer()
}

// Records a card move into a GameState
func (g *Game) RecordMove(cardID string, inputs []interface{}) error {
	// Find the card
	cardIndex := -1
	for i, c := range g.State.Cards {
		if c.ID == cardID {
			cardIndex = i
			break
		}
	}
	if cardIndex == -1 {
		return fmt.Errorf("card with ID %s not found", cardID)
	}

	// Validade input length
	expected := CountRequiredInputs(g.State.Cards[cardIndex].InputsReq)
	if len(inputs) != expected {
		return fmt.Errorf("expected %d inputs but got %d", expected, len(inputs))
	}

	// Store inputs in GameState's Card
	card := &g.State.Cards[cardIndex]
	card.Inputs = inputs

	// Add the cardID to the queue
	g.State.Queue = append(g.State.Queue, cardID)

	return nil
}

// Count the number of required inputs, aka lowercase letters
func CountRequiredInputs(req string) int {
	count := 0
	for _, ch := range req {
		if ch >= 'a' && ch <= 'z' {
			count++
		}
	}
	return count
}
