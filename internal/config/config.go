package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr        string
	BaseHTTP    string
	LogLevel    string
	FileName    string
	DatabaseDSN string
	JWTSecret   string
}

const (
	defaultAddr      = "localhost:8080"
	defaultBaseHTTP  = "http://localhost:8080"
	defaultLogLevel  = "debug"
	defaultJWTSecret = "im_not_a_secret_plz_replace_me"
)

func NewConfig() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", defaultAddr, "Адрес сервера в формате <хост>:<порт>")
	flag.StringVar(&config.BaseHTTP, "b", defaultBaseHTTP, "HTTP адрес сервера в сокращенном URL в формате <http схема>://<хост>:<порт>")
	flag.StringVar(&config.LogLevel, "l", defaultLogLevel, "Уровень логирования")
	flag.StringVar(&config.FileName, "f", "", "JSON файл с данными")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Адрес подключения к базе данных")
	flag.StringVar(&config.JWTSecret, "t", defaultJWTSecret, "Секретный ключ JWT")

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

	databaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		config.DatabaseDSN = databaseDSN
	}

	return &config
}
