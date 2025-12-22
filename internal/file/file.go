package file

import (
	"encoding/json"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/dustin/go-humanize"
	"os"
)

type (
	ShorterItem struct {
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}
	ShorterItems []ShorterItem
)

func NewShorterItems() *ShorterItems {
	return &ShorterItems{}
}

// LoadShorterItems получает данные из файла и десериализует массив байтов в объект `ShorterItems`
func (items *ShorterItems) LoadShorterItems(filename string) (*ShorterItems, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return items, err
	}
	if len(data) == 0 {
		return items, fmt.Errorf("файл пуст")
	}

	err = json.Unmarshal(data, items)
	if err != nil {
		return items, err
	}
	logger.Log.Debugf("Найдено записей в файле %s: %d", filename, len(*items))

	return items, nil
}

// SaveShorterItems сохраняет данные из storage в filename
func (items *ShorterItems) SaveShorterItems(filename string, storage repository.Repository) error {
	logger.Log.Debugf("Сохраняем записи storage в файл %s", filename)
	for id, originURL := range storage.GetAll() {
		*items = append(*items, ShorterItem{ShortUrl: id, OriginalUrl: originURL})
	}
	itemsNum := len(*items)
	logger.Log.Debugf("Найдено записей в storage: %d", itemsNum)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	itemsSize, err := file.Write(data)
	if err != nil {
		return err
	}

	logger.Log.Debugf("В файл сохранено %d записей(%s)", itemsNum, humanize.Bytes(uint64(itemsSize)))
	defer file.Close()

	return nil
}
