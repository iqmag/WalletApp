package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"WalletApp/internal/domain"
)

const (
	DEPOSIT  = "DEPOSIT" // депозит
	WITHDRAW = "WITHDRAW" // снятие
)

// WalletService бизнес-логика кошельков
type WalletService interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
	CreateWallet(ctx context.Context) (uuid.UUID, error)
	PerformOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int64) error
}

// Структура walletService реализует интерфейс WalletService
type walletService struct {
	repo domain.WalletRepository
}

// экземпляр
func NewWalletService(repo domain.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

// Метод для получения текущего баланса кошелька по id
func (s *walletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(ctx, walletID)
}

// Метод для создания нового кошелька и возврата его id
func (s *walletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	return s.repo.CreateWallet(ctx)
}

// Метод для выполнения операций (депозит/снятие) с кошельком
func (s *walletService) PerformOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int64) error {
	switch operationType {
	case DEPOSIT:
		return s.repo.UpdateBalance(ctx, walletID, amount) // Увеличиваем баланс
	case WITHDRAW:
		balance, err := s.repo.GetBalance(ctx, walletID)
		if err != nil {
			return err
		}
		if balance < amount {
			return errors.New("insufficient funds") // Ошибка при недостатке средств
		}
		return s.repo.UpdateBalance(ctx, walletID, -amount) // Уменьшаем баланс
	default:
		return errors.New("invalid operation type") // Ошибка при неверном типе операции
	}
}