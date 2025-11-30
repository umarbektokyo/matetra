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
}
