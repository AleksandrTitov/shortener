package id

import "math/rand"

const (
	LenID   = 6
	Symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GetID() string {
	id := make([]byte, LenID)

	for i := range len(id) {
		id[i] = Symbols[rand.Intn(len(Symbols))]
	}
	return string(id)
}
