package repository

import (
	"context"

	"github.com/plusik10/cmd/order-info-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	GetOrderUIDs(ctx context.Context) ([]string, error)
	GetOrderByUID(ctx context.Context, orderUID string) (model.Order, error)
}
