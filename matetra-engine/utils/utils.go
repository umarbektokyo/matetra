package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var VERSION = "0.1"
var PORT int = 1729
var DECK_PATH string = "cards/cards.csv"

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
