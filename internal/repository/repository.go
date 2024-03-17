package repository

import (
	"context"

	"github.com/plusik10/cmd/order-info-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
}
