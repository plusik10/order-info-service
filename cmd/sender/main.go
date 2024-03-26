package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/plusik10/cmd/order-info-service/internal/config"
	"github.com/plusik10/cmd/order-info-service/internal/model"
)

func nConnect(path string) (*nats.Conn, error) {
	for {
		time.Sleep(1 * time.Second)
		nc, err := nats.Connect(path)
		if err == nil {
			return nc, nil
		} else {
			log.Println(err.Error())
			log.Println("Reconnecting to NATS...")
		}
	}
}
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := nConnect(cfg.Nuts.URL) // TODO: add config file
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(cfg.Nuts.ClusterID, cfg.Nuts.ClientPubId, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	log.Println("Connecting to nuts successfully")
	i := 0
	for {
		time.Sleep(5 * time.Second)
		i = i + 1
		if i%5 == 0 {
			jsonData, err := randomJson()
			if err != nil {
				log.Println("Error generating json data err: ", err.Error())
				continue
			}

			err = sc.Publish(cfg.Nuts.Subject, jsonData)
			if err != nil {
				log.Println("Error publishing to NATS err: ", err.Error())
				continue
			}

			log.Println("Successfully published badjson data")
		} else {
			OrderUID := gofakeit.UUID()
			cost := gofakeit.Uint16()
			goodsTotal := gofakeit.Uint16()
			payment := model.Payment{
				Transaction:  gofakeit.UUID(),
				RequestID:    gofakeit.UUID(),
				Currency:     gofakeit.CurrencyShort(),
				Provider:     gofakeit.CurrencyLong(),
				Amount:       int(cost + goodsTotal),
				PaymentDt:    int(gofakeit.Uint32()),
				Bank:         gofakeit.BeerName(),
				DeliveryCost: int(cost),
				GoodsTotal:   int(goodsTotal),
				CustomFee:    int(gofakeit.Uint16()),
			}

			delivery := model.Delivery{
				Name:    gofakeit.Name(),
				Phone:   gofakeit.Phone(),
				Zip:     gofakeit.Zip(),
				City:    gofakeit.City(),
				Address: gofakeit.CountryAbr(),
				Region:  gofakeit.CurrencyShort(),
				Email:   gofakeit.Email(),
			}

			countItem := gofakeit.Uint32() % 10
			Items := make([]model.Item, 0)

			for i := 0; i < int(countItem); i++ {
				Items = append(Items, model.Item{
					ChrtID:      int(gofakeit.Uint8()),
					TrackNumber: gofakeit.UUID(),
					Price:       int(gofakeit.Uint8()),
					Rid:         gofakeit.UUID(),
					Name:        gofakeit.Name(),
					Sale:        int(gofakeit.Uint8()),
					Size:        "0",
					TotalPrice:  int(gofakeit.Uint8()),
					NmID:        int(gofakeit.Uint8()),
					Brand:       gofakeit.Company(),
					Status:      gofakeit.StatusCode(),
				})
			}

			order := model.Order{
				OrderUID:          OrderUID,
				TrackNumber:       gofakeit.StreetNumber(),
				Entry:             "DHL",
				Delivery:          delivery,
				Payment:           payment,
				Items:             Items,
				Locale:            "en-US",
				InternalSignature: gofakeit.UUID(),
				CustomerID:        gofakeit.UUID(),
				DeliveryService:   gofakeit.UUID(),
				Shardkey:          gofakeit.UUID(),
				SmID:              1,
				DateCreated:       gofakeit.Date(),
				OofShard:          gofakeit.UUID(),
			}

			jsonData, err := json.Marshal(order)
			if err != nil {
				log.Println("Error generating json data err: ", err.Error())
				continue
			}

			err = sc.Publish(cfg.Nuts.Subject, jsonData)
			if err != nil {
				log.Println("Error publishing to NATS err: ", err.Error())
				continue
			}
			log.Println("order publish success: ", OrderUID)
		}
	}
}

type badJsonData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func randomJson() ([]byte, error) {
	djd := badJsonData{
		Name: gofakeit.Name(),
		Age:  int(gofakeit.Int16()),
	}

	return json.Marshal(djd)
}
