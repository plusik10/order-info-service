package app

import (
	"context"
	"log"

	"github.com/plusik10/cmd/order-info-service/internal/config"
	"github.com/plusik10/cmd/order-info-service/internal/repository"
	"github.com/plusik10/cmd/order-info-service/internal/repository/postgres"
	"github.com/plusik10/cmd/order-info-service/pkg/db"
)

type serviceProvider struct {
	config          *config.Config
	db              db.Client
	orderRepository repository.OrderRepository
	//orderService    service.OrderService
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		config: cfg,
	}
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
