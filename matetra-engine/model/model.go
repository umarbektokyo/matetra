package model

type Card struct {
	ID          string // Unique Identifier for every card, even if a dublicate id is different
	Name        string
	Description string
	Type        string
	Method      string        // Defines what method in code will be taken
	Owner       int           // -1: deck, -2: used, User.ID: owner
	Inputs      []interface{} // list of ints and str
	InputsReq   string        // string with each character signifying input number type.
	// InputsReq explained:
	// d: dice (int)
	// O: makes the next digit optional
	// U: makes next digit user's
	// A: Makes next digit attacked one
	// n: number (int)
	// c: card, id (string)
	// p: player (int)
	// X: minimum for the input (int)
	// Y: maximum for the input (int)
	// i: allow input for a user (int)
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
