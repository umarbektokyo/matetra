package cards

import (
	"encoding/csv"
	"fmt"
	"matetra/model"
	"matetra/utils"
	"os"
	"strconv"
)

// Loads cards from csv file
func LoadCardsFromCSV(path string) ([]model.Card, error) {
	// Open the file
	file := utils.Must(os.Open("cards.csv"))
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
		count := utils.Must(strconv.Atoi(row[7]))

		// Add multiple copies if necessary
		for i := 0; i < count; i++ {
			idCounter++
			card := model.Card{
				ID:          fmt.Sprintf("%s_%d", row[0], idCounter),
				Name:        row[0],
				Description: row[1],
				Type:        row[2],
				Method:      row[3],
				InputsReq:   row[4],
				Owner:       -1,
				Target:      [][2]int{},
				Inputs:      [][2]int{},
			}
			cards = append(cards, card)
		}
	}
	return cards, nil
}
