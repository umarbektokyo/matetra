package cards

import (
	"encoding/csv"
	"fmt"
	"matetra/cards/functions"
	"matetra/model"
	"matetra/utils"
	"os"
	"strconv"
)

// Loads cards from csv file
func LoadCardsFromCSV(path string) ([]model.Card, error) {
	// Open the file
	file := utils.Must(os.Open(path))
	defer file.Close()

	// Read through the file and clean the data
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records := utils.Must(reader.ReadAll())
	records = records[1:]

	// Start recording the cards
	var cards []model.Card

	// Add each card (row)
	for _, row := range records {
		count := utils.Must(strconv.Atoi(row[6]))

		// Add multiple copies if necessary
		for i := 0; i < count; i++ {
			card := model.Card{
				Name:        row[0],
				Description: row[2],
				Type:        row[3],
				Method:      row[4],
				InputsReq:   row[5],
				Owner:       -1,
				Inputs:      []int{},
			}
			cards = append(cards, card)
		}
	}
	return cards, nil
}

func CardFunction(vgs *model.GameState, cardIndex int) error {
	var card *model.Card

	card = &vgs.Cards[cardIndex]

	switch card.Method {
	case "ADD":
		return functions.ADD(vgs, card)

	case "SUBTRACT":
		return functions.SUBTRACT(vgs, card)

	default:
		return fmt.Errorf("unknown card method %s", card.Method)
	}
}
