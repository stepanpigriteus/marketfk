package live

import (
	"context"
	"log"
	"sync"
	"time"

	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
)

func GenConnectAndRead(port string, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()

	// для теста отключения
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	defer cancel()

	address := "exchange" + port[4:5] + ":" + port
	exchange := model.Exchange{Name: "Exchange" + port[4:5]}

	log.Printf("Connecting to %s", address)

	client := live.NewLiveExchangeClient(exchange, address)

	err := client.Connect()
	if err != nil {
		log.Printf("Error connecting to %s: %v", address, err)
		return
	}

	priceCh := make(chan model.Price, 100)

	client.StartReading(ctx, priceCh)

	for price := range priceCh {
		output <- price // отправка в общий input
	}
	log.Printf("Stopped reading from %s", address)

	// повторное подключение
}
