package live

import (
	"fmt"
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
	"sync"
)

type Quote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func ConnectAndRead(port string, wg *sync.WaitGroup) {
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

	priceCh := make(chan model.Price)

	client.StartReading(priceCh)

	for price := range priceCh {
		fmt.Printf("Port %s: %+v\n", port, price)
	}

	log.Printf("Stopped reading from %s", address)
}
