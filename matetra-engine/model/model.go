package model

type Card struct {
	ID          string // Unique Identifier for every card, even if a dublicate id is different
	Name        string
	Description string
	Type        string
	Method      string   // Defines what method in code will be taken
	Owner       int      // -1: deck, -2: used, User.ID: owner
	Target      [][2]int // list of [Player.id, Numbers[index]] for every target
	Inputs      [][2]int // list of [Player.ID, Numbers[index]]
	InputsReq   string   // string with each character signifying input number type.
	// InputsReq explained:
	// d: dice
	// O: makes the next digit optional
	// U: makes next digit user's
	// A: Makes next digit attacked one
	// n: number
	// c: card
	// p: player
}

// Only for authentication
type Player struct {
	Name string
	Hash string
}

// Main Game Object
type GameState struct {
	GameID  string
	Players []Player
	Cards   []Card
	Numbers [][5]string
	Done    []bool
	Queue   []string // stores CardID's and every time the move is finished, we apply all the cards and cleane the data in them, marking them as used.
	Turn    int      // total turns elapsed; current player = Turn % len(Players)
}
