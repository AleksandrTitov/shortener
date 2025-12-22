package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr     string
	BaseHTTP string
	LogLevel string
	FileName string
}

const (
	defaultAddr     = "localhost:8080"
	defaultBaseHTTP = "http://localhost:8080"
	defaultLogLevel = "info"
)

func NewConfig() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", defaultAddr, "Адрес сервера в формате <хост>:<порт>")
	flag.StringVar(&config.BaseHTTP, "b", defaultBaseHTTP, "HTTP адрес сервера в сокращенном URL в формате <http схема>://<хост>:<порт>")
	flag.StringVar(&config.LogLevel, "l", defaultLogLevel, "Уровень логирования")
	flag.StringVar(&config.FileName, "f", "", "JSON файл с данными")

	flag.Parse()

	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		config.Addr = addr
	}

	baseHTTP, ok := os.LookupEnv("BASE_URL")
	if ok {
		config.BaseHTTP = baseHTTP
	}

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		config.LogLevel = logLevel
	}

	fileName, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileName = fileName
	}

	return &config
}
