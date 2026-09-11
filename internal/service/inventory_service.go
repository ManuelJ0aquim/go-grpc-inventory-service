package service

import (
	"context"

	"github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/domain"
)

type InventoryService interface {
	CheckStock(ctx context.Context, productID string, quantity int32) (bool, int32, error)
	ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, string, error)
}

type inventoryService struct {
	repo domain.Repository
}

func NewInventoryService(repo domain.Repository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CheckStock(ctx context.Context, productID string, quantity int32) (bool, int32, error) {
	stock, err := s.repo.GetStock(productID)
	if err != nil {
		return false, 0, err
	}
	return stock >= quantity, stock, nil
}

func (s *inventoryService) ReserveStock(ctx context.Context, orderID, productID string, quantity int32) (bool, string, error) {
	err := s.repo.Reserve(productID, quantity)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "Estoque reservado com sucesso", nil
}