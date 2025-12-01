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
	idCounter := 0

	// Add each card (row)
	for _, row := range records {
		count := utils.Must(strconv.Atoi(row[6]))

		// Add multiple copies if necessary
		for i := 0; i < count; i++ {
			idCounter++
			card := model.Card{
				ID:          fmt.Sprintf("%s_%d", row[0], idCounter),
				Name:        row[0],
				Description: row[2],
				Type:        row[3],
				Method:      row[4],
				InputsReq:   row[5],
				Owner:       -1,
				Inputs:      []interface{}{},
			}
			cards = append(cards, card)
		}
	}
	return cards, nil
}

func CardFunction(vgs *model.GameState, cardID string) {
	// Get the currently used card
	var card *model.Card
	for i := range vgs.Cards {
		if vgs.Cards[i].ID == cardID {
			card = &vgs.Cards[i]
			break
		}
	}

	if card == nil {
		return
	}

	switch card.Method {
	case "ADD":
		functions.ADD(vgs, card)
	case "SUBTRACT":
		functions.SUBTRACT(vgs, card)
	}
}
