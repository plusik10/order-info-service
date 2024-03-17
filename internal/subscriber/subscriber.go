package subscriber

import (
	"context"
	"encoding/json"

	"github.com/nats-io/stan.go"
	"github.com/plusik10/cmd/order-info-service/internal/model"
	"github.com/plusik10/cmd/order-info-service/internal/service"
)

type Subscriber struct {
	Nc           stan.Conn
	sub          stan.Subscription
	orderService service.OrderService
}

func NewSubscriber(clusterID string, clientID string, service service.OrderService) (*Subscriber, error) {
	nc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}
	return &Subscriber{Nc: nc, orderService: service}, nil
}

func (s *Subscriber) Create(ctx context.Context, data []byte) error {
	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}
	return s.orderService.Create(ctx, order)
}

func (s *Subscriber) Subscribe(subject string, callback stan.MsgHandler, opts ...stan.SubscriptionOption) error {
	sub, err := s.Nc.Subscribe(subject, callback, opts...)
	if err != nil {
		return err
	}
	s.sub = sub
	return nil
}

func (s *Subscriber) Unsubscribe() {
	if s.sub != nil {
		s.sub.Unsubscribe()
	}
}

func (s *Subscriber) Close() {
	if s.Nc != nil {
		s.Nc.Close()
	}
}
