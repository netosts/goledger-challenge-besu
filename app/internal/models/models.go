package models

import (
	"errors"
	"time"
)

type StoredValue struct {
	ID        int       `json:"id" db:"id"`
	Value     uint64    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SetValueRequest struct {
	Value uint64 `json:"value" binding:"required"`
}

type ValueResponse struct {
	Value uint64 `json:"value"`
}

type CheckResponse struct {
	IsEqual         bool   `json:"is_equal"`
	DatabaseValue   uint64 `json:"database_value"`
	BlockchainValue uint64 `json:"blockchain_value"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r *SetValueRequest) IsValid() error {
	if r.Value > 1e18 {
		return errors.New("value is too large")
	}

	return nil
}
