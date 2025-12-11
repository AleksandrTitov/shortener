package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Addr     string
	BaseHTTP string
}

func NewConfig() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", "localhost:8080", "Адрес сервера в формате <хост>:<порт>")
	flag.StringVar(&config.BaseHTTP, "b", "http://localhost:8080", "HTTP адрес сервера в сокращенном URL в формате <http схема>://<хост>:<порт>")

	flag.Parse()

	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		config.Addr = addr
	}

	baseHTTP, ok := os.LookupEnv("BASE_URL")
	if ok {
		config.BaseHTTP = baseHTTP
	}

	log.Printf("INFO: Адрес сервера %s", config.Addr)

	return &config
}
