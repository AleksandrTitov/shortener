package config

import (
	"flag"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"os"
)

type Config struct {
	Addr     string
	BaseHTTP string
}

const (
	defaultAddr     = "localhost:8080"
	defaultBaseHTTP = "http://localhost:8080"
)

func NewConfig() *Config {
	var config Config

	log := logger.NewLogger()

	flag.StringVar(&config.Addr, "a", defaultAddr, "Адрес сервера в формате <хост>:<порт>")
	flag.StringVar(&config.BaseHTTP, "b", defaultBaseHTTP, "HTTP адрес сервера в сокращенном URL в формате <http схема>://<хост>:<порт>")

	flag.Parse()

	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		config.Addr = addr
	}

	baseHTTP, ok := os.LookupEnv("BASE_URL")
	if ok {
		config.BaseHTTP = baseHTTP
	}

	log.Infof("Адрес сервера %s", config.Addr)

	return &config
}
