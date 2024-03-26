package app

import (
	"context"
	"log"

	"github.com/plusik10/cmd/order-info-service/internal/cache"
	localcach "github.com/plusik10/cmd/order-info-service/internal/cache/localCach"
	"github.com/plusik10/cmd/order-info-service/internal/config"
	"github.com/plusik10/cmd/order-info-service/internal/repository"
	"github.com/plusik10/cmd/order-info-service/internal/repository/postgres"
	"github.com/plusik10/cmd/order-info-service/internal/service"
	"github.com/plusik10/cmd/order-info-service/internal/service/order"
	"github.com/plusik10/cmd/order-info-service/pkg/db"
)

type serviceProvider struct {
	config          *config.Config
	db              db.Client
	orderRepository repository.OrderRepository
	orderService    service.OrderService
	cache           cache.Cache
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	sp := &serviceProvider{config: cfg}
	return sp
}

func (s *serviceProvider) GetConfig() *config.Config {
	if s.config == nil {
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("failed to create config: %s", err.Error())
		}
		s.config = cfg
	}
	return s.config
}

func (s *serviceProvider) GetDB(ctx context.Context) db.Client {
	if s.db == nil {
		cfg, err := s.GetConfig().GetDBConfig()
		if err != nil {
			log.Fatalf("failed to get db config: %s", err.Error())
		}

		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			log.Fatalf("can't connect to db: %v", err.Error)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetOrderRepository(ctx context.Context) repository.OrderRepository {
	if s.orderRepository == nil {
		s.orderRepository = postgres.NewOrderRepository(s.GetDB(ctx))
	}

	return s.orderRepository
}

func (s *serviceProvider) GetCache() cache.Cache {
	if s.cache == nil {
		s.cache = localcach.NewLocalCache(s.GetConfig().DefaultExpiration, s.GetConfig().CleanupInterval)
	}

	return s.cache
}

func (s *serviceProvider) GetOrderService(ctx context.Context) service.OrderService {
	if s.orderService == nil {
		s.orderService = order.NewOrderService(s.GetOrderRepository(ctx), s.GetCache())
		err := s.orderService.LoadOrders(ctx)
		if err != nil {
			log.Println("Error loading order service")
		}

	}

	return s.orderService
}
