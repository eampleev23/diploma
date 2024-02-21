package cnf

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	RanAddr     string
	LogLevel    string
	DBDSN       string
	SecretKey   string
	TLimitQuery time.Duration
	TokenExp    time.Duration
}

func NewConfig() (*Config, error) {
	config := &Config{
		TLimitQuery: 20 * time.Second, //nolint:gomnd //20 секунд максимальнео время на запрос к бд
		TokenExp:    time.Hour * 1344, //nolint:gomnd //2 месяца не истекает авторизация
	}
	err := config.SetValues()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) SetValues() error {
	// регистрируем переменную flagRunAddr как аргумент -a со значением по умолчанию :8080
	flag.StringVar(&c.RanAddr, "a", "localhost:8080", "address and port to run server")
	// регистрируем уровень логирования
	flag.StringVar(&c.LogLevel, "l", "debug", "logger level")
	// принимаем строку подключения к базе данных
	flag.StringVar(&c.DBDSN, "d", "", "postgres database")
	// принимаем секретный ключ сервера для авторизации
	flag.StringVar(&c.SecretKey, "s", "e4853f5c4810101e88f1898db21c15d3", "server's secret key for authorization")
	// парсим переданные серверу аргументы в зарегестрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		c.RanAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.LogLevel = envLogLevel
	}

	if envDBDSN := os.Getenv("DATABASE_DSN"); envDBDSN != "" {
		c.DBDSN = envDBDSN
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		c.SecretKey = envSecretKey
	}
	return nil
}
