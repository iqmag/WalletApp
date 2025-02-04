package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// структура для хранения конфигурации приложения
type Config struct {
    DBUrl string // URL бд
}

// LoadConfig загружает конфигурацию из файла .env или из переменных окружения
func LoadConfig() *Config {
    err := godotenv.Load() // передаем env файл
    if err != nil {
        log.Println("Warning: .env file not found or cannot be loaded")
    }

    return &Config{
        DBUrl: os.Getenv("DATABASE_URL"),
    }
}