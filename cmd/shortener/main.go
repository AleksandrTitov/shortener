package main

import (
	"github.com/AleksandrTitov/shortener/internal/app"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"os"
)

func main() {
	a, err := InitializeApp()
	if err != nil {
		logger.Log.Fatal(err.Error())
		os.Exit(2)
	}

	if err = a.Run(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func InitializeApp() (*app.App, error) {
	err := logger.Initialize("info", true)
	if err != nil {
		logger.Log.Warnf("Не удалось инициализировать логгер: %v", err.Error())
	}
	a := app.NewApp()

	return a, nil
}
