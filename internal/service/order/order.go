package order

import (
	"context"
	"encoding/json"
	"errors"
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

func (s *orderService) Create(ctx context.Context, data []byte) error {
	var order model.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}
	if err := order.Validate(); err != nil {
		return fmt.Errorf("unvalid order json err: %v", err)
	}
	s.Cache.Set(order.OrderUID, data, 0)
	return s.Repo.Create(ctx, order)
}

// GetOrderUIDs Возвращает список uid заказов
func (s *orderService) GetOrderUIDs(ctx context.Context) ([]string, error) {
	return s.Repo.GetOrderUIDs(ctx)
}

// GetOrderByUID returns order by uid
func (s *orderService) GetOrderByUID(ctx context.Context, orderUID string) ([]byte, error) {
	o, found := s.Cache.Get(orderUID)
	if !found {
		order, err := s.Repo.GetOrderByUID(ctx, orderUID)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(order)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	switch v := o.(type) {
	case nil:
		return nil, errors.New("invalid order")
	case []byte:
		return v, nil
	default:
		return nil, errors.New("invalid order")
	}
}

func (s *orderService) LoadOrders(ctx context.Context) error {
	orders, err := s.Repo.GetOrders(ctx)
	if err != nil {
		return err
	}

	for _, o := range orders {
		order, err := json.Marshal(o)
		if err != nil {
			return err
		}
		s.Cache.Set(o.OrderUID, order, 0)
	}
	return nil
}
