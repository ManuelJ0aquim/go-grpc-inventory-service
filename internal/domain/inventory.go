package domain

import "errors"

var (
	ErrProductNotFound   = errors.New("produto não encontrado")
	ErrInsufficientStock = errors.New("estoque insuficiente")
)

type Item struct {
	ProductID string
	Quantity  int32
}

type Repository interface {
	GetStock(productID string) (int32, error)
	Reserve(productID string, quantity int32) error
}