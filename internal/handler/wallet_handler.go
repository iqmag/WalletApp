package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"WalletApp/internal/usecase"
)

// структура WalletHandler для обработки запросов, связанных с кошельками
type WalletHandler struct {
	Service usecase.WalletService // Сервис для выполнения операций с кошельками
	Logger  *logrus.Logger // логгер

	Semaphore chan struct{} // Семафор для ограничения параллельных запросов
}

// экземпляр
func NewWalletHandler(service usecase.WalletService, logger *logrus.Logger) *WalletHandler {
	return &WalletHandler{
		Service:   service,
		Logger:    logger,
		Semaphore: make(chan struct{}, 10), // Ограничиваем до 10 параллельных запросов
	}
}

// Метод для обработки создания нового кошелька
func (h *WalletHandler) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	walletID, err := h.Service.CreateWallet(ctx) // Вызов метода create
	if err != nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError) // Возврат ошибки 500 при неудаче
		return
	}

	response := map[string]uuid.UUID{"walletId": walletID}
	w.Header().Set("Content-Type", "application/json") // заголовок на JSON
	w.WriteHeader(http.StatusCreated) // код 201 (Создано)
	json.NewEncoder(w).Encode(response) // Кодирование ответа в JSON и отправка клиенту
}

// Метод для обработки запроса на получение баланса кошелька
func (h *WalletHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	walletIDParam := mux.Vars(r)["walletId"] // Извлечение id кошелька из URL
	walletID, err := uuid.Parse(walletIDParam) // Парсинг id кошелька из строки

	if err != nil || walletID == uuid.Nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest) // Возврат ошибки 400 при некорректном id
		return
	}

	ctx := r.Context()
	balance, err := h.Service.GetBalance(ctx, walletID)

	if err != nil {
		http.Error(w, "Error retrieving balance", http.StatusInternalServerError) // Возврат ошибки 500 при неудаче
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
}

// Метод для обработки операций (депозит/снятие) с кошельком
func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	h.Semaphore <- struct{}{} // Захватываем слот для выполнения запроса
	defer func() { <-h.Semaphore }() // Освобождаем слот после завершения

	var request struct {
		WalletId      uuid.UUID `json:"walletId"`
		OperationType string    `json:"operationType"`
		Amount        int64     `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	err := h.Service.PerformOperation(r.Context(), request.WalletId, request.OperationType, request.Amount)

	if err != nil {
        if err.Error() == "insufficient funds" {
            http.Error(w, "Insufficient funds", http.StatusBadRequest) // Возврат ошибки 400 при недостатке средств
            return
        }
        if err.Error() == "invalid operation type" { 
            http.Error(w, "Invalid operation type", http.StatusBadRequest) // Возврат ошибки 400 при неверном типе операции
            return
        }
        http.Error(w, "Error performing operation", http.StatusInternalServerError) // Возврат ошибки 500 при неудаче выполнения операции
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}