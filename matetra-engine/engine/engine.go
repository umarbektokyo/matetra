package engine

import (
	"errors"
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
			Turn:    0,
		},
	}
}

// Return the index of the player whoose turn i is
func (g *Game) CurrentPlayer() int {
	n := len(g.State.Players)
	if n == 0 {
		return -1
	}
	return g.State.Turn % n
}

// Adds a new player to the game
func (g *Game) AddPlayer(name, hash string) int {
	p := model.Player{
		Name:    name,
		Hash:    hash,
		Numbers: [5]string{},
		Hand:    [6]string{},
		Queue:   []string{},
		Done:    false,
	}
	g.State.Players = append(g.State.Players, p)
	return len(g.State.Players) - 1
}

// Adds a new card to the game
func (g *Game) AddCard(card model.Card) {
	g.State.Cards = append(g.State.Cards, card)
}

// Adds a card to the player's hand
func (g *Game) GiveCardToPlayer(cardID string, playerID int) error {
	if playerID < 0 || playerID >= len(g.State.Players) {
		return errors.New("invalid playerID")
	}

	for i := range g.State.Cards {
		if g.State.Cards[i].ID == cardID {
			for h := range g.State.Players[playerID].Hand {
				if g.State.Players[playerID].Hand[h] == "" {
					g.State.Players[playerID].Hand[h] = cardID
					g.State.Cards[i].Owner = playerID
					return nil
				}
			}
			return errors.New("hand is full")
		}
	}
	return errors.New("card not found")
}

// Move a card into an action // we should figure out a better solution

// Toggles Done boolean for the player
func (g *Game) ToggleTurn(playerID int) error {
	if playerID < 0 || playerID >= len(g.State.Players) {
		return errors.New("invalid playerID")
	}
	g.State.Players[playerID].Done = !g.State.Players[playerID].Done
	return nil
}

// Checks if everyone ended the turn
func (g *Game) TurnEnded() bool {
	for i := range g.State.Players {
		if !g.State.Players[i].Done {
			return false
		}
	}
	return true
}

func (g *Game) EndTurn() error {
	if !g.TurnEnded() {
		return errors.New("some people still haven't finished the turn")
	}
	// Reset Done variables, fill hands and increment turn.
	for i := range g.State.Players {
		g.State.Players[i].Done = false
		// CONTINUE FROM HERE
		// - update number values
		// - reset used cards
		// - fill players' hands
		// - whatever else I've missed
	}
	g.State.Turn += 1

	return nil
}
