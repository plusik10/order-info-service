package service

import (
	"context"
)

type OrderService interface {
	Create(ctx context.Context, data []byte) error
	GetOrderUIDs(ctx context.Context) ([]string, error)
	GetOrderByUID(ctx context.Context, orderUID string) ([]byte, error)
	LoadOrders(ctx context.Context) error
}
