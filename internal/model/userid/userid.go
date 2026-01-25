package userid

import (
	"fmt"
	"github.com/google/uuid"
)

func New() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("ошибка генерации user ID: %v", err)
	}
	return id.String(), nil
}
