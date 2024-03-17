package main

import (
	"context"
	"log"
	"os"

	"github.com/plusik10/cmd/order-info-service/internal/app"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx, os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
