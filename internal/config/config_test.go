package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// временная переменная окружения
	os.Setenv("DATABASE_URL", "postgres://postgres:cyball@localhost:5432/postgres?sslmode=disable")
	defer os.Unsetenv("DATABASE_URL")

	cfg := LoadConfig()

	// Проверяем, что значение DBUrl соответствует ожидаемому
	expectedDBUrl := "postgres://postgres:cyball@localhost:5432/postgres?sslmode=disable"
	if cfg.DBUrl != expectedDBUrl {
		t.Errorf("expected DATABASE_URL to be %s, got %s", expectedDBUrl, cfg.DBUrl)
	}
}