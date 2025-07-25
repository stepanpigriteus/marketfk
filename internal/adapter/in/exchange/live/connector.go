package live

import (
	"context"
	"fmt"
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
	"sync"
	"time"
)

func GenConnectAndRead(ctx context.Context, host string, port string, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()
	var address string
	exchange := model.Exchange{}
	if host == "" {
		address = "127.0.0.1:" + port
		exchange = model.Exchange{Name: "ExchangeT" + port[4:5]}
	} else {
		address = host + port[4:5] + ":" + port
		exchange = model.Exchange{Name: "Exchange" + port[4:5]}
	}

	fmt.Println(address)

	const maxRetries = 5
	const maxBackoff = 16 * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Printf("Завершение подключения к %s по сигналу", address)
			return
		default:
		}

		log.Printf("Подключение к %s", address)
		client := live.NewLiveExchangeClient(exchange, address)

		err := client.Connect()
		if err != nil {
			log.Printf("Ошибка подключения к %s: %v", address, err)
			if !tryReconnect(ctx, client, address, maxRetries, maxBackoff) {
				log.Printf("Не удалось переподключиться к %s после %d попыток", address, maxRetries)
				return
			}
			continue
		}

		priceCh := make(chan model.Price, 100)

		// Старт чтения
		client.StartReading(ctx, priceCh) // важно: передаём ctx

	READ_LOOP:
		for {
			select {
			case price, ok := <-priceCh:
				if !ok {
					break READ_LOOP
				}
				select {
				case output <- price:
				case <-ctx.Done():
					break READ_LOOP
				}
			case <-ctx.Done():
				break READ_LOOP
			}
		}

		log.Printf("Соединение с %s разорвано, попытка переподключения", address)
		if !tryReconnect(ctx, client, address, maxRetries, maxBackoff) {
			log.Printf("Не удалось переподключиться к %s после %d попыток", address, maxRetries)
			return
		}
	}
}

func tryReconnect(ctx context.Context, client *live.LiveExchangeClient, address string, maxRetries int, maxBackoff time.Duration) bool {
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			log.Printf("Переподключение к %s отменено", address)
			return false
		default:
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			log.Printf("Попытка переподключения %d/%d к %s через %v", attempt+1, maxRetries, address, backoff)
			time.Sleep(backoff)

			err := client.Connect()
			if err == nil {
				log.Printf("Успешное переподключение к %s", address)
				return true
			}
			log.Printf("Ошибка переподключения к %s: %v", address, err)
		}
	}
	return false
}
