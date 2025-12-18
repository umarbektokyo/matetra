package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"matetra/model"
	"math/rand"
	"time"
)

var VERSION = "0.1"
var PORT int = 1729
var DECK_PATH string = "cards/cards.csv"
var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func MatetraSplash() {
	fmt.Println("Matetra v" + VERSION)
}

func ValidateInputs(vgs *model.GameState, card *model.Card) error {
	// Check the length
	if len(card.Inputs) != len(card.InputsReq) {
		return fmt.Errorf(
			"%s expects %d inputs, got %d",
			card.Method, len(card.InputsReq), len(card.Inputs),
		)
	}

	// Check input values
	for i := range card.InputsReq {
		val := card.Inputs[i]
		switch card.InputsReq[i] {
		case 'd':
			if val < 1 || val > 6 {
				return fmt.Errorf("input %d must be dice (1..6), got %v", i, val)
			}

		case 'p':
			if val < 0 || val >= len(vgs.Players) {
				return fmt.Errorf("input %d must be player index, got %v", i, val)
			}

		case 'U':
			if val < 0 || val >= len(vgs.Players) {
				return fmt.Errorf("input %d must be player index, got %v", i, val)
			}
			if val != card.Owner {
				return fmt.Errorf("input %d must be your own index (%v), got %v", i, card.Owner, val)
			}

		case 'A':
			if val < 0 || val >= len(vgs.Players) {
				return fmt.Errorf("input %d must be player index, got %v", i, val)
			}

			if val != (vgs.Turn % len(vgs.Players)) {
				return fmt.Errorf("input %d must be index of defending player, got %v", i, val)
			}

		case 'n':
			if val < 0 || val > 4 {
				return fmt.Errorf("input %d must be number index 0..4, got %v", i, val)
			}

		case 'c':
			fmt.Println("You forgot to implement this!!")

		case 'i':
			X := card.Inputs[i-2]
			Y := card.Inputs[i-1]
			if val < X || val > Y {
				return fmt.Errorf("input %d must be in range of %d..%d, got %d", i, X, Y, val)
			}
		}

	}

	return nil
}

func RollDice(sides int) int {
	roll := r.Intn(sides) + 1
	return roll
}
