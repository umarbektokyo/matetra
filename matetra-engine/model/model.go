package model

type Card struct {
	ID          string // Unique Identifier for every card, even if a dublicate id is different
	Name        string
	Description string
	Type        string
	Method      string   // Defines what method in code will be taken
	Owner       int      // -1: deck, -2: used, User.ID: owner
	Target      [2]int   // [Player.id, Numbers[index]]
	Inputs      [][2]int // list of [Player.ID, Numbers[index]]
}

type Player struct {
	Name    string
	Hash    string
	Numbers [5]string
	Hand    [6]string
	Queue   []string // Card.ID
	Done    bool
}

type GameState struct {
	GameID  string
	Players []Player
	Cards   []Card
	Turn    int // total turns elapsed; current player = Turn % len(Players)
}
