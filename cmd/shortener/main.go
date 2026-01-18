package main

import (
	"github.com/AleksandrTitov/shortener/internal/app"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"os"
)

func main() {
	a, err := InitializeApp()

	if a.DB != nil {
		defer a.DB.Close()
	}
	if err != nil {
		logger.Log.Fatal(err.Error())
		os.Exit(2)
	}

	logger.Log.Fatal(a.Run())
}

func InitializeApp() (*app.App, error) {
	err := logger.Initialize("info", true)
	if err != nil {
		logger.Log.Warnf("Не удалось инициализировать логгер: %v", err.Error())
	}
	a := app.NewApp()

	return a, nil
}
