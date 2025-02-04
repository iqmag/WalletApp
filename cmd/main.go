package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	 _ "github.com/golang-migrate/migrate/v4/source/file"

	"WalletApp/internal/config"
	"WalletApp/internal/handler"
	"WalletApp/internal/repository"
	"WalletApp/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Настройка миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Путь к директории с миграциями
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	repo := repository.NewPostgresWalletRepository(db)
	service := usecase.NewWalletService(repo)
	logger := logrus.New()

	h := handler.NewWalletHandler(service, logger)

	r := mux.NewRouter()

	// Добавляем обработчик для корневого маршрута
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to WalletApp API"))
	})

	// Регистрация маршрутов API
	r.HandleFunc("/api/v1/wallet", h.HandleCreateWallet).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/wallets/{walletId}", h.HandleGetBalance).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/wallets/operation", h.HandleOperation).Methods(http.MethodPost)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}