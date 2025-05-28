package live

import (
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
	"sync"
)

func GenConnectAndRead(port string, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()

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
	client.StartReading(priceCh)

	for price := range priceCh {
		output <- price // отправка в общий input
	}
	log.Printf("Stopped reading from %s", address)
}
