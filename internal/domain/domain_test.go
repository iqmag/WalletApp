package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestWalletCreation(t *testing.T) {
	// Создаем новый кошелек
	walletID := uuid.New()
	initialBalance := int64(1000)
	wallet := Wallet{
		ID:      walletID,
		Balance: initialBalance,
	}

	// Проверяем, что ID установлен правильно
	if wallet.ID != walletID {
		t.Errorf("expected wallet ID %s, got %s", walletID.String(), wallet.ID.String())
	}

	// Проверяем, что баланс установлен правильно
	if wallet.Balance != initialBalance {
		t.Errorf("expected balance %d, got %d", initialBalance, wallet.Balance)
	}
}