package engine

import (
	"fmt"
	"matetra/cards"
	"matetra/model"
	"matetra/utils"
	"math/big"
	"math/rand"
	"sync"
)

type Game struct {
	State *model.GameState
	mu    sync.RWMutex
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
			Queue:   make([]int, 0),
			Turn:    -1,
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
func (g *Game) AddPlayer(name, hash string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Adds a new player object
	playerID := len(g.State.Players)
	g.State.Players = append(g.State.Players, model.Player{
		Name: name,
		Hash: hash,
	})
	g.State.Numbers = append(g.State.Numbers, NewNumberRow())
	g.State.Done = append(g.State.Done, false)
	return playerID, nil
}

// Return the index of the player whoose turn it is
func (g *Game) CurrentPlayer() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	n := len(g.State.Players)
	if n == 0 {
		return -1
	}
	return g.State.Turn % n
}

// Loads the card deck into game state
func (g *Game) LoadCards() {
	g.mu.Lock()
	defer g.mu.Unlock()

	deck := utils.Must(cards.LoadCardsFromCSV(utils.DECK_PATH))
	g.State.Cards = append(g.State.Cards, deck...)
}

// Check if everyone has finished the turn
func (g *Game) TurnsFinished() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, done := range g.State.Done {
		if !done {
			return false
		}
	}
	return true
}

// Checks how many cards does the player have
func (g *Game) PlayerHandCount(player int) int {
	g.mu.RLock()
	defer g.mu.RUnlock()

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
	g.mu.Lock()
	defer g.mu.Unlock()

	for p := range g.State.Players {
		handCount := 0
		for _, card := range g.State.Cards {
			if card.Owner == p {
				handCount++
			}
		}

		for handCount < 6 {
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

			if len(deck) == 0 {
				break
			}
			// Choose a card from a deck
			idx := deck[rand.Intn(len(deck))]
			g.State.Cards[idx].Owner = p
			handCount++
		}
	}
}

// Makes a virtual deep copy of the game state
func (g *Game) CopyState() *model.GameState {
	g.mu.RLock()
	defer g.mu.RUnlock()

	virtual := &model.GameState{
		GameID:  g.State.GameID,
		Players: append([]model.Player(nil), g.State.Players...),
		Cards:   append([]model.Card(nil), g.State.Cards...),
		Numbers: make([][5]model.Number, len(g.State.Numbers)),
		Done:    append([]bool(nil), g.State.Done...),
		Queue:   append([]int(nil), g.State.Queue...),
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
func (g *Game) ApplyCard(vgs *model.GameState, cardIndex int) error {
	err := cards.CardFunction(vgs, cardIndex)
	if err != nil {
		return err
	}

	// Remove card after applying
	vgs.Cards[cardIndex].Owner = -2
	vgs.Cards[cardIndex].Inputs = nil

	return nil
}

// Applies all the Cards in Queue
func (g *Game) ApplyCards(vgs *model.GameState) error {
	for _, cardIndex := range vgs.Queue {
		err := g.ApplyCard(vgs, cardIndex)
		if err != nil {
			return err
		}
	}
	vgs.Queue = nil

	return nil
}

// Ends a player's turn
func (g *Game) NextTurn() error {
	virtual := g.CopyState()
	if err := g.ApplyCards(virtual); err != nil {
		return err
	}
	g.State = virtual
	g.RestockCards()
	g.State.Turn++
	for i := range g.State.Done {
		g.State.Done[i] = false
	}
	return nil
}

// Checks if the player is moving their own card
func (g *Game) PlayerCanPlayCard(playerID, cardIndex int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if cardIndex < 0 || cardIndex >= len(g.State.Cards) {
		return false
	}

	return g.State.Cards[cardIndex].Owner == playerID
}

// API: Moves
func (g *Game) ProcessMove(playerID int, cardIndex int, inputs []int, permanent bool) (*model.GameState, error) {
	g.mu.RLock()

	// check ownership
	if !g.PlayerCanPlayCard(playerID, cardIndex) {
		return nil, fmt.Errorf("you do not own this card")
	}

	// validate input
	expected := len(g.State.Cards[cardIndex].InputsReq)
	if len(inputs) != expected {
		return nil, fmt.Errorf("expected %d inputs but got %d", expected, len(inputs))
	}

	g.mu.RUnlock()

	if !permanent {
		virtual := g.CopyState()

		vCard := &virtual.Cards[cardIndex]
		vCard.Inputs = append([]int(nil), inputs...)

		if err := g.ApplyCard(virtual, cardIndex); err != nil {
			return nil, fmt.Errorf("preview failed: %v", err)
		}

		return virtual, nil
	} else {
		g.mu.Lock()
		defer g.mu.Unlock()

		if g.State.Done[playerID] {
			return nil, fmt.Errorf("you have already finished your turn")
		}

		g.State.Cards[cardIndex].Inputs = append([]int(nil), inputs...)
		g.State.Queue = append(g.State.Queue, cardIndex)

		return g.State, nil
	}
}

// API: Turns
func (g *Game) ProcessNextTurn(playerID int) (*model.GameState, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State.Done[playerID] {
		return nil, fmt.Errorf("you have already finished your turn")
	}

	g.State.Done[playerID] = true

	finished := true
	for _, done := range g.State.Done {
		if !done {
			finished = false
			break
		}
	}

	if finished {
		if err := g.NextTurn(); err != nil {
			return nil, err
		}
	}

	return g.State, nil
}
