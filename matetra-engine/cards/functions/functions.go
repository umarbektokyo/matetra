package functions

import (
	"fmt"
	"matetra/model"
	"math/big"
)

// Input: AnUn
func ADD(vgs *model.GameState, card *model.Card) error {
	ValidateInputs(vgs, card)

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
	ValidateInputs(vgs, card)

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
	ValidateInputs(vgs, card)

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
	ValidateInputs(vgs, card)

	attackerPlayer := card.Inputs[0]
	attackerIndex := card.Inputs[1]
	userPlayer := card.Inputs[2]
	userIndex := card.Inputs[3]

	a := &vgs.Numbers[attackerPlayer][attackerIndex]
	b := &vgs.Numbers[userPlayer][userIndex]

	a.Value.Div(a.Value, b.Value)

	b.Value = big.NewInt(0)
	b.Mark = "n"

	return nil
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

		case 'p', 'U', 'A':
			if val < 0 || val > len(vgs.Players) {
				return fmt.Errorf("input %d must be player index, got %v", i, val)
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
