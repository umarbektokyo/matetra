package functions

import (
	"fmt"
	"matetra/model"
	"matetra/utils"
	"math/big"
)

// Input: AnUn
func ADD(vgs *model.GameState, card *model.Card) error {
	utils.ValidateInputs(vgs, card)

	attackerPlayer := card.Inputs[0]
	attackerIndex := card.Inputs[1]
	userPlayer := card.Inputs[2]
	userIndex := card.Inputs[3]

	a := &vgs.Numbers[attackerPlayer][attackerIndex]
	b := &vgs.Numbers[userPlayer][userIndex]

	a.Value.Add(a.Value, b.Value)

	b.Value = big.NewInt(0)
	b.Mark = "n"

	return nil
}

// Input: AnUn
func SUBTRACT(vgs *model.GameState, card *model.Card) error {
	utils.ValidateInputs(vgs, card)

	attackerPlayer := card.Inputs[0]
	attackerIndex := card.Inputs[1]
	userPlayer := card.Inputs[2]
	userIndex := card.Inputs[3]

	a := &vgs.Numbers[attackerPlayer][attackerIndex]
	b := &vgs.Numbers[userPlayer][userIndex]

	a.Value.Sub(a.Value, b.Value)

	b.Value = big.NewInt(0)
	b.Mark = "n"

	return nil
}

// Input: AnUn
func MULTIPLY(vgs *model.GameState, card *model.Card) error {
	utils.ValidateInputs(vgs, card)

	attackerPlayer := card.Inputs[0]
	attackerIndex := card.Inputs[1]
	userPlayer := card.Inputs[2]
	userIndex := card.Inputs[3]

	a := &vgs.Numbers[attackerPlayer][attackerIndex]
	b := &vgs.Numbers[userPlayer][userIndex]

	a.Value.Mul(a.Value, b.Value)

	b.Value = big.NewInt(0)
	b.Mark = "n"

	return nil
}

// Input: AnUn
func DIVIDE(vgs *model.GameState, card *model.Card) error {
	utils.ValidateInputs(vgs, card)

	attackerPlayer := card.Inputs[0]
	attackerIndex := card.Inputs[1]
	userPlayer := card.Inputs[2]
	userIndex := card.Inputs[3]

	a := &vgs.Numbers[attackerPlayer][attackerIndex]
	b := &vgs.Numbers[userPlayer][userIndex]

	if b.Value == big.NewInt(0) {
		return fmt.Errorf("Cannot divide by zero")
	}

	a.Value.Div(a.Value, b.Value)

	b.Value = big.NewInt(0)
	b.Mark = "n"

	return nil
}
