package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"WalletApp/internal/handler"
)

type mockWalletService struct {
	balances map[uuid.UUID]int64 // балансы кошельков по id
}

// экземпляр
func newMockWalletService() *mockWalletService {
	return &mockWalletService{
		balances: make(map[uuid.UUID]int64),
	}
}

// Метод для выполнения операций (депозит/снятие) с кошельком
func (m *mockWalletService) PerformOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int64) error {
	if operationType == "DEPOSIT" { // // Проверка на (депозит)
		m.balances[walletID] += amount
		return nil
	}
	if operationType == "WITHDRAW" { // Проверка на (снятие)
		if m.balances[walletID] < amount {
			return errors.New("insufficient funds")
		}
		m.balances[walletID] -= amount
		return nil
	}
	return errors.New("invalid operation type")
}

// Метод для получения текущего баланса кошелька
func (m *mockWalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	balance, exists := m.balances[walletID]
	if !exists {
		return 0, errors.New("wallet not found")
	}
	return balance, nil
}

// Метод для создания нового кошелька и возврата его id
func (m *mockWalletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	newID := uuid.New()
	m.balances[newID] = 0 // Начальный баланс = 0
	return newID, nil
}

// Тестирование функции обработки операций с кошельком
func TestHandleOperation(t *testing.T) {
	logger := logrus.New()
	mockSvc := newMockWalletService() // экземпляр имитированного сервиса кошельков
	h := handler.NewWalletHandler(mockSvc, logger)

	ctx := context.Background()
	walletID, _ := mockSvc.CreateWallet(ctx) // Создаем новый кошелёк

	tests := []struct {
		name           string
		operationType  string
		amount         int64
		expectedCode   int
		expectedBalance int64
	}{
		{"Deposit valid amount", "DEPOSIT", 500, http.StatusOK, 500}, // Тест на успешный депозит
		{"Withdraw valid amount", "WITHDRAW", 200, http.StatusOK, 300}, // Тест на успешное снятие
		{"Withdraw more than balance", "WITHDRAW", 1000, http.StatusBadRequest, 500}, // Тест на снятие больше чем баланс
		{"Invalid operation type", "INVALID", 500, http.StatusBadRequest, 500}, // Тест на неверный тип операции
	}

	for _, tt := range tests {
	    t.Run(tt.name, func(t *testing.T) { // Запуск под-теста с именем текущего случая
	        body := map[string]interface{}{
	            "walletId":      walletID.String(), // id
	            "operationType": tt.operationType, // тип
	            "amount":        tt.amount, // сумма
	        }
	        jsonBody, _ := json.Marshal(body)
	        req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets/operation", bytes.NewBuffer(jsonBody))
	        w := httptest.NewRecorder()

	        h.HandleOperation(w, req)

	        res := w.Result()
	        if res.StatusCode != tt.expectedCode {
	            t.Errorf("expected status %d, got %d", tt.expectedCode, res.StatusCode)
	        }

	        if res.StatusCode == http.StatusOK { // Если операция успешна
	            balanceAfterOperation, _ := mockSvc.GetBalance(ctx, walletID) // получаем баланс после операции
	            if balanceAfterOperation != tt.expectedBalance { // сравниваем полученный баланс с ожидаемым значением
	                t.Errorf("expected balance after operation %d, got %d", tt.expectedBalance, balanceAfterOperation) // ошибка при несоответствии балансов
	            }
	        }
	    })
    }
}