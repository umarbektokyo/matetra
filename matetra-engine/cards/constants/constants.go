package constants

import (
	"fmt"
	"matetra/model"
	"math/big"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func DICE(vgs *model.GameState, player int) error {

	foundSlot := -1
	for i, num := range vgs.Numbers[player] {
		if num.Mark == "n" {
			foundSlot = i
			break
		}
	}

	if foundSlot == -1 {
		return fmt.Errorf("no empty slots available to roll a dice")
	}

	roll := r.Intn(6) + 1

	vgs.Numbers[player][foundSlot].Value = big.NewFloat(float64(roll))
	vgs.Numbers[player][foundSlot].Mark = ""

	return nil
}

func CONSTPI(vgs *model.GameState, card *model.Card) error {
	userPlayer := card.Inputs[0]
	userIndex := card.Inputs[1]

	a := &vgs.Numbers[userPlayer][userIndex]

	a.Value = nil
	a.Mark = ""

	return nil
}
