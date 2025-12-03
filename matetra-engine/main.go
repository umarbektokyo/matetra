package main

import (
	"fmt"
	"matetra/engine"
)

func main() {
	game := engine.New("8f67b270e1768b444ac717b07a626529fa22487a78feec3281ea8c67ccc74235")
	game.LoadCards()
	fmt.Println("Development in progress...")
	fmt.Println("Current GameState")
	fmt.Println(game)
	game.AddPlayer("Umarbek", "d80827dcc407c5382f3656f3e2d5488a0a52b1f1eb27059da515c77f4a8fec88")
	game.AddPlayer("Cassini", "e3d5951820c5c01627af6b60a49cf47d2ae8ed37e8b7db020e5e93e992421c40")
	game.NextTurn()
	// game.RecordMove(car)
}
