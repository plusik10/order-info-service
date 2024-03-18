package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/stan.go"
	"github.com/plusik10/cmd/order-info-service/internal/config"
	"github.com/plusik10/cmd/order-info-service/internal/subscriber"
)

type App struct {
	serviceProvider *serviceProvider
	pathConfig      string
	subscriber      *subscriber.Subscriber
}

func NewApp(ctx context.Context, pathConfig string) (*App, error) {
	a := &App{pathConfig: pathConfig}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, err
}

func (a *App) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	a.StartSubscriber(ctx, &wg)
	wg.Wait()
	return nil
}

func (a *App) StartSubscriber(ctx context.Context, wg *sync.WaitGroup) error {
	subject := a.serviceProvider.GetConfig().Subject
	a.subscriber.Nc.Subscribe(subject, func(msg *stan.Msg) {
		err := a.subscriber.Create(ctx, msg.Data)
		if err != nil {
			log.Println("Error creating order err: ", err.Error())
		}
		log.Println("Order created successfully")
	})
	<-ctx.Done()
	wg.Done()
	a.subscriber.Close()
	a.subscriber.Unsubscribe()
	fmt.Println("end subscriber")
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		a.initSubscribe,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
func (a *App) initSubscribe(ctx context.Context) error {
	clusterID := a.serviceProvider.GetConfig().Nuts.ClusterID
	clientID := a.serviceProvider.GetConfig().Nuts.ClientSubID

	s, err := subscriber.NewSubscriber(clusterID, clientID, a.serviceProvider.GetOrderService(ctx))
	if err != nil {
		return err
	}
	a.subscriber = s
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	op := "initServiceProvider"
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf(op+": %s", err.Error())
	}

	a.serviceProvider = newServiceProvider(cfg)

	return nil
}
