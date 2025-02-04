package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"WalletApp/internal/handler"
)

// MockWalletService представляет собой мок-сервис для тестирования
type MockWalletService struct {
	Balances map[uuid.UUID]int64
}

func NewMockWalletService() *MockWalletService {
	return &MockWalletService{
		Balances: make(map[uuid.UUID]int64),
	}
}

func (m *MockWalletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	walletID := uuid.New()
	m.Balances[walletID] = 0
	return walletID, nil
}

func (m *MockWalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	balance, exists := m.Balances[walletID]
	if !exists {
		return 0, errors.New("wallet not found")
	}
	return balance, nil
}

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

// Тест корневого маршрута "/"
func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to WalletApp API"))
	})

	r.ServeHTTP(w, req)

	res := w.Result()
	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "Welcome to WalletApp API", body.String())
}

// Тест создания кошелька
func TestCreateWallet(t *testing.T) {
	service := NewMockWalletService()
	h := handler.NewWalletHandler(service, logrus.New())

	reqBody := map[string]interface{}{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.HandleCreateWallet(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotEmpty(t, response["walletId"]) // Проверяем наличие ID кошелька в ответе
}

// Тест получения баланса кошелька
func TestGetBalance(t *testing.T) {
	service := NewMockWalletService()
	h := handler.NewWalletHandler(service, logrus.New())

	walletID, _ := service.CreateWallet(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
	wr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/{walletId}", h.HandleGetBalance).Methods(http.MethodGet)
	r.ServeHTTP(wr, req)

	res := wr.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]int64
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, int64(0), response["balance"]) // Проверяем начальный баланс
}

// Тест выполнения операции с кошельком (депозит)
func TestPerformOperation_Deposit(t *testing.T) {
	service := NewMockWalletService()
	h := handler.NewWalletHandler(service, logrus.New())

	walletID, _ := service.CreateWallet(context.Background())

	reqBody := map[string]interface{}{
		"walletId":      walletID.String(),
		"operationType": "DEPOSIT",
		"amount":        100,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets/operation", bytes.NewBuffer(body))
	wr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/operation", h.HandleOperation).Methods(http.MethodPost)
	r.ServeHTTP(wr, req)

	res := wr.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	balanceResp := service.Balances[walletID]
	assert.Equal(t, int64(100), balanceResp) // Проверяем баланс после депозита
}

// Тест выполнения операции с кошельком (вывод средств)
func TestPerformOperation_Withdraw(t *testing.T) {
	service := NewMockWalletService()
	h := handler.NewWalletHandler(service, logrus.New())

	walletID, _ := service.CreateWallet(context.Background())
	service.PerformOperation(context.Background(), walletID, "DEPOSIT", 100)

	reqBody := map[string]interface{}{
		"walletId":      walletID.String(),
		"operationType": "WITHDRAW",
		"amount":        50,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets/operation", bytes.NewBuffer(body))
	wr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/operation", h.HandleOperation).Methods(http.MethodPost)
	r.ServeHTTP(wr, req)

	res := wr.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	balanceResp := service.Balances[walletID]
	assert.Equal(t, int64(50), balanceResp) // Проверяем баланс после вывода средств
}