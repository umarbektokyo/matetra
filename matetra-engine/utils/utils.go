package utils

var PORT int = 1729
var DECK_PATH string = "cards/cards.csv"

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
