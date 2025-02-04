package domain

import (
	"context"

	"github.com/google/uuid"
)

// интерфейс для работы с кошельками
type WalletRepository interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
	CreateWallet(ctx context.Context) (uuid.UUID, error)
}