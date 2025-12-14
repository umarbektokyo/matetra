package constants

import (
	"matetra/model"
	"matetra/utils"
)

func CONSTPI(vgs *model.GameState, card *model.Card) error {
	utils.ValidateInputs(vgs, card)

	userPlayer := card.Inputs[0]
	userIndex := card.Inputs[1]

	a := &vgs.Numbers[userPlayer][userIndex]

	a.Value = nil
	a.Mark = "PI"

	return nil
}
