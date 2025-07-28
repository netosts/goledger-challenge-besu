package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/netosts/goledger-challenge-besu/internal/models"
)

type Repository interface {
	GetValue() (*models.StoredValue, error)
	SetValue(value uint64) error
	Close() error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetValue() (*models.StoredValue, error) {
	var sv models.StoredValue
	query := `SELECT id, value, created_at, updated_at FROM stored_values ORDER BY updated_at DESC LIMIT 1`

	err := r.db.QueryRow(query).Scan(&sv.ID, &sv.Value, &sv.CreatedAt, &sv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get value: %w", err)
	}

	return &sv, nil
}

func (r *PostgresRepository) SetValue(value uint64) error {
	query := `
	UPDATE stored_values 
	SET value = $1, updated_at = $2 
	WHERE id = (SELECT id FROM stored_values ORDER BY updated_at DESC LIMIT 1)
	`

	_, err := r.db.Exec(query, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
