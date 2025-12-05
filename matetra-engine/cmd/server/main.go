package main

import (
	"fmt"
	"log"
	"matetra/api"
	"matetra/engine"
	"matetra/utils"
	"os"
)

func main() {
	cmd := os.Args

	if len(cmd) == 1 {
		utils.MatetraSplash()
		return
	}

	switch cmd[1] {
	case "start":
		if len(cmd) < 3 {
			fmt.Println("please provide a title for the game.")
			fmt.Println("(ex: matetra start wonderful-game)")
			return
		}

		game := engine.New(cmd[2])

		log.Printf("loading card deck...")
		game.LoadCards()
		log.Printf("deck loaded with %d cards", len(game.State.Cards))

		apiServer := api.New(game)
		apiServer.Start()
	default:
		utils.MatetraSplash()
		fmt.Println(cmd[1] + " not recognised.")
	}
}
