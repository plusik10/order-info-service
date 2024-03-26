package service

import (
	"context"

	"github.com/plusik10/cmd/order-info-service/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, order model.Order) error
	GetOrderUIDs(ctx context.Context) ([]string, error)
	GetOrderByUID(ctx context.Context, orderUID string) (model.Order, error)
	LoadOrders(ctx context.Context) error
}
