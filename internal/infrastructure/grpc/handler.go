package grpc

import (
	"context"

	inventoryv1 "github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/infrastructure/grpc/pb/inventory/v1"
	"github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	svc service.InventoryService
}

func NewHandler(svc service.InventoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CheckStock(ctx context.Context, req *inventoryv1.CheckStockRequest) (*inventoryv1.CheckStockResponse, error) {
	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id e obrigatorio")
	}

	available, currentStock, err := h.svc.CheckStock(ctx, req.GetProductId(), req.GetQuantity())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "erro ao consultar estoque: %v", err)
	}

	return &inventoryv1.CheckStockResponse{
		Available:    available,
		CurrentStock: currentStock,
	}, nil
}

func (h *Handler) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	if req.GetProductId() == "" || req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id e product_id sao obrigatorios")
	}

	if req.GetQuantity() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantidade deve ser maior que zero")
	}

	success, message, err := h.svc.ReserveStock(ctx, req.GetOrderId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		return &inventoryv1.ReserveStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &inventoryv1.ReserveStockResponse{
		Success: success,
		Message: message,
	}, nil
}