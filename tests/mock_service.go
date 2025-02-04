package tests

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// MockWalletService мок-сервис для тестирования
type MockWalletService struct {
	Balances map[uuid.UUID]int64
}

// экземпляр
func NewMockWalletService() *MockWalletService {
	return &MockWalletService{
		Balances: make(map[uuid.UUID]int64),
	}
}

// метод для создания нового кошелька и возврата id
func (m *MockWalletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	walletID := uuid.New()
	m.Balances[walletID] = 0
	return walletID, nil
}

// метод для получения текущего баланса кошелька по id
func (m *MockWalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	balance, exists := m.Balances[walletID]
	if !exists {
		return 0, errors.New("wallet not found")
	}
	return balance, nil
}

// метод для выполнения операций (депозит/снятие) с кошельком
func (m *MockWalletService) PerformOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int64) error {
	switch operationType {
	case "DEPOSIT":
		m.Balances[walletID] += amount
	case "WITHDRAW":
		if m.Balances[walletID] < amount {
			return errors.New("insufficient funds")
		}
		m.Balances[walletID] -= amount
	default:
		return errors.New("invalid operation type")
	}
	return nil
}