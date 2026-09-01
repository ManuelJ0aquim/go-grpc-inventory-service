package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetStock(productID string) (int32, error) {
	query := `SELECT quantity FROM inventory WHERE product_id = $1`
	var quantity int32

	err := r.db.QueryRow(query, productID).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrProductNotFound
		}
		return 0, fmt.Errorf("erro ao buscar estoque: %w", err)
	}

	return quantity, nil
}

func (r *PostgresRepository) Reserve(productID string, quantity int32) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transacao: %w", err)
	}
	defer tx.Rollback()

	// Lock pessimista na linha para evitar overselling
	var currentStock int32
	querySelect := `SELECT quantity FROM inventory WHERE product_id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, querySelect, productID).Scan(&currentStock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("erro ao consultar produto com lock: %w", err)
	}

	if currentStock < quantity {
		return domain.ErrInsufficientStock
	}

	queryUpdate := `UPDATE inventory SET quantity = quantity - $1, updated_at = NOW() WHERE product_id = $2`
	_, err = tx.ExecContext(ctx, queryUpdate, quantity, productID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estoque: %w", err)
	}

	return tx.Commit()
}