package concurrency

import (
	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
)

func FanIn(clients []*live.LiveExchangeClient) <-chan model.Price {
	out := make(chan model.Price)

	for _, client := range clients {
		client.StartReading(out)
	}

	return out
}
