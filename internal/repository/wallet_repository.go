package repository

import (
    "context"
    "database/sql"

    "github.com/google/uuid"
)

// структура PostgresWalletRepository для работы с кошельками в Postgres
type PostgresWalletRepository struct {
    db *sql.DB
}

// экземпляр
func NewPostgresWalletRepository(db *sql.DB) *PostgresWalletRepository {
    return &PostgresWalletRepository{db: db}
}

// Метод для получения баланса кошелька по id
func (r *PostgresWalletRepository) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
    var balance int64
    err := r.db.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1", walletID).Scan(&balance)
    if err == sql.ErrNoRows {
        return 0, nil
    }
    return balance, err
}

// Метод для обновления баланса кошелька по id
func (r *PostgresWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
    _, err := r.db.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
    return err
}

// Метод для создания нового кошелька и возврата id
func (r *PostgresWalletRepository) CreateWallet(ctx context.Context) (uuid.UUID, error) {
    walletID := uuid.New()
    _, err := r.db.ExecContext(ctx, "INSERT INTO wallets (id, balance) VALUES ($1, $2)", walletID, 0)
    return walletID, err
}