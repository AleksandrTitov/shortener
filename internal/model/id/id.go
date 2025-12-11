package id

import (
	"errors"
	"math/rand"
)

const (
	LenID   = 6
	Symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	GetIDError = errors.New("ошибка генерации")
)

type GeneratorID interface {
	GetID() (string, error)
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (*Generator) GetID() (string, error) {
	id := make([]byte, LenID)

	for i := range len(id) {
		id[i] = Symbols[rand.Intn(len(Symbols))]
	}
	return string(id), nil
}
