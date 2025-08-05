package usecases

import "github.com/netosts/goledger-challenge-besu/internal/models"

// ContractUseCaseInterface defines the interface for contract use case operations
type ContractUseCaseInterface interface {
	SetValue(value uint64) error
	GetValue() (uint64, error)
	SyncValue() error
	CheckValue() (*models.CheckResponse, error)
	Close()
}
