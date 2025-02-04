package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"WalletApp/internal/usecase"
)

type mockWalletRepository struct {
	wallets map[uuid.UUID]int64 // Хранилище для кошельков
}

func newMockWalletRepository() *mockWalletRepository {
	return &mockWalletRepository{
		wallets: make(map[uuid.UUID]int64),
	}
}

func (m *mockWalletRepository) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	balance, exists := m.wallets[walletID]
	if !exists {
		return 0, errors.New("wallet not found")
	}
	return balance, nil
}

func (m *mockWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	if _, exists := m.wallets[walletID]; !exists {
		return errors.New("wallet not found")
	}
	m.wallets[walletID] += amount // Обновляем баланс
	return nil
}

func (m *mockWalletRepository) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	walletID := uuid.New()
	m.wallets[walletID] = 0 // Создаём кошелёк с балансом 0
	return walletID, nil
}

func TestPerformOperation(t *testing.T) {
	repo := newMockWalletRepository()
	svc := usecase.NewWalletService(repo) // Передаём мок-репозиторий

	tests := []struct {
		name           string
		operationType  string
		amount         int64
		expectError    bool
		expectedBalance int64 // для проверки баланса после операции
	}{
		{"Deposit valid amount", usecase.DEPOSIT, 500, false, 500},
		{"Withdraw valid amount", usecase.WITHDRAW, 200, false, 300},
		{"Withdraw more than balance", usecase.WITHDRAW, 1000, true, 0}, // Ожидаем ошибку
        {"Invalid operation type", "INVALID", 500, true, -1},            // Ожидаем ошибку
    }

	for _, tt := range tests {
	    t.Run(tt.name, func(t *testing.T) {
	        walletID, _ := repo.CreateWallet(context.Background()) // Создаём новый кошелёк
	        repo.UpdateBalance(context.Background(), walletID, 0) // Устанавливаем начальный баланс в 0

	        if tt.operationType == usecase.DEPOSIT {

	        // Начальный баланс остается 0, депозит добавит 500
	        } else if tt.operationType == usecase.WITHDRAW {
	            repo.UpdateBalance(context.Background(), walletID, 500); // Устанавливаем начальный баланс в 500 для тестов вывода
	        }

	        err := svc.PerformOperation(context.Background(), walletID ,tt.operationType ,tt.amount)

	        if (err != nil) != tt.expectError {
	            t.Errorf("expected error status %v but got %v", tt.expectError, err)
	        }

	        // Проверяем баланс после операции только если операция успешная
	        if !tt.expectError {
	            balanceAfterOperation, _ := repo.GetBalance(context.Background(), walletID)
	            if balanceAfterOperation != tt.expectedBalance {
	                t.Errorf("expected balance after operation %d, got %d", tt.expectedBalance, balanceAfterOperation)
	            }
	        }
	    })
    }
}