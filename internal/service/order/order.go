package order

import (
	"context"
	"fmt"

	"github.com/plusik10/cmd/order-info-service/internal/cache"
	"github.com/plusik10/cmd/order-info-service/internal/model"
	"github.com/plusik10/cmd/order-info-service/internal/repository"
	"github.com/plusik10/cmd/order-info-service/internal/service"
)

type orderService struct {
	Repo  repository.OrderRepository
	Cache cache.Cache
}

func NewOrderService(repo repository.OrderRepository, cache cache.Cache) service.OrderService {
	return &orderService{
		Repo:  repo,
		Cache: cache,
	}
}

func (s *orderService) Create(ctx context.Context, order model.Order) error {
	// validate order
	s.Cache.Set(order.OrderUID, order, 0)
	return s.Repo.Create(ctx, order)
}

// GetOrderUIDs Возвращает список uid заказов
func (s *orderService) GetOrderUIDs(ctx context.Context) ([]string, error) {
	return s.Repo.GetOrderUIDs(ctx)
}

// GetOrderByUID returns order by uid
func (s *orderService) GetOrderByUID(ctx context.Context, orderUID string) (model.Order, error) {
	var order model.Order
	o, found := s.Cache.Get(orderUID)
	if !found {
		return s.Repo.GetOrderByUID(ctx, orderUID)
	}
	orderPtr, ok := o.(*model.Order)
	if !ok {
		return order, fmt.Errorf("failed to convert order to *model.Order")
	}

	return *orderPtr, nil
}

func (s *orderService) LoadOrders(ctx context.Context) error {

	orders, err := s.Repo.GetOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.Cache.Set(order.OrderUID, order, 0)
	}
	return nil
}
