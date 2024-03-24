package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/nats-io/stan.go"
	"github.com/plusik10/cmd/order-info-service/internal/api/v1/handlers/order"
	"github.com/plusik10/cmd/order-info-service/internal/config"
	"github.com/plusik10/cmd/order-info-service/internal/subscriber"
)

type App struct {
	serviceProvider *serviceProvider
	pathConfig      string
	httpServer      *http.Server
	subscriber      *subscriber.Subscriber
}

func NewApp(ctx context.Context, pathConfig string) (*App, error) {
	a := &App{pathConfig: pathConfig}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.serviceProvider.db.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	go a.StartSubscriber(ctx, wg)
	go func() {
		err := a.runPublicHTTP()
		if err != nil {
			log.Println("error running: ", err)
		}
		wg.Done()
	}()
	wg.Wait()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	}
	log.Println("HTTP server shutdown")

	return nil
}

func (a *App) initPublicHttp(ctx context.Context) error {
	r := chi.NewRouter()
	r.Get("/", order.GetOrderUIDs(ctx, a.serviceProvider.GetOrderService(ctx)))
	r.Get("/order/info", order.Info(ctx, a.serviceProvider.GetOrderService(ctx)))
	a.httpServer = &http.Server{
		Handler:           r,
		Addr:              a.serviceProvider.GetConfig().HTTP.Port,
		ReadHeaderTimeout: a.serviceProvider.GetConfig().HTTP.Timeout,
	}

	return nil
}

func (a *App) runPublicHTTP() error {
	fmt.Println("Starting public http server")
	if err := a.httpServer.ListenAndServe(); err != nil {
		return err
	}

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
	wg.Done()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		a.initPublicHttp,
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
