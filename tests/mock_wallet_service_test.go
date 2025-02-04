package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	service := NewMockWalletService()
	ctx := context.Background()

	walletID, err := service.CreateWallet(ctx)
	assert.NoError(t, err) // Проверяем, что ошибки нет
	assert.NotEqual(t, uuid.Nil, walletID) // Проверяем, что ID кошелька не пустой
	assert.Equal(t, int64(0), service.Balances[walletID]) // Проверяем, что баланс равен 0
}

func TestGetBalance(t *testing.T) {
	service := NewMockWalletService()
	ctx := context.Background()

	// Создаем кошелек
	walletID, _ := service.CreateWallet(ctx)

	// Проверяем баланс существующего кошелька
	balance, err := service.GetBalance(ctx, walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), balance)

	// Проверяем ошибку для несуществующего кошелька
	_, err = service.GetBalance(ctx, uuid.New())
	assert.Error(t, err)
	assert.Equal(t, "wallet not found", err.Error())
}

func TestPerformOperation_Deposit(t *testing.T) {
	service := NewMockWalletService()
	ctx := context.Background()

	walletID, _ := service.CreateWallet(ctx)

	// Депозит 100
	err := service.PerformOperation(ctx, walletID, "DEPOSIT", 100)
	assert.NoError(t, err)

	balance, _ := service.GetBalance(ctx, walletID)
	assert.Equal(t, int64(100), balance) // Баланс должен быть 100
}

func TestPerformOperation_Withdraw(t *testing.T) {
	service := NewMockWalletService()
	ctx := context.Background()

	walletID, _ := service.CreateWallet(ctx)

	// Депозит 100
	service.PerformOperation(ctx, walletID, "DEPOSIT", 100)

	// Снятие 50
	err := service.PerformOperation(ctx, walletID, "WITHDRAW", 50)
	assert.NoError(t, err)

	balance, _ := service.GetBalance(ctx, walletID)
	assert.Equal(t, int64(50), balance) // Баланс должен быть 50

	// Попытка снять больше средств
	err = service.PerformOperation(ctx, walletID, "WITHDRAW", 100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}

func TestPerformOperation_InvalidOperationType(t *testing.T) {
	service := NewMockWalletService()
	ctx := context.Background()

	walletID, _ := service.CreateWallet(ctx)

	// Неверный тип операции
	err := service.PerformOperation(ctx, walletID, "INVALID_OP", 10)
	assert.Error(t, err)
	assert.Equal(t, "invalid operation type", err.Error())
}