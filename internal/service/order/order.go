package order

import (
	"context"

	"github.com/plusik10/cmd/order-info-service/internal/model"
	"github.com/plusik10/cmd/order-info-service/internal/repository"
	"github.com/plusik10/cmd/order-info-service/internal/service"
)

type orderService struct {
	Repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) service.OrderService {
	return &orderService{
		Repo: repo,
	}
}

func (s *orderService) Create(ctx context.Context, order model.Order) error {
	return s.Repo.Create(ctx, order)
}
