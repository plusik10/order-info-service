package app

import (
	"context"

	"github.com/plusik10/cmd/order-info-service/internal/config"
)

type App struct {
	serviceProvider *serviceProvider
	pathConfig      string
	//subscriber      *subscriber.Subscriber
}

func NewApp(ctx context.Context, pathConfig string) *App {
	a := &App{pathConfig: pathConfig}
	err := a.initDeps(ctx)
	if err != nil {
		return nil
	}

	return a
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		//a.initSubscribe,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	a.serviceProvider = newServiceProvider(cfg)

	return nil
}
