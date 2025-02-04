package domain

import "github.com/google/uuid"

// модель
type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}