package live

import (
	"context"
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/exchange/live"
	"marketfuck/internal/domain/model"
	"sync"
	"time"
)

func GenConnectAndRead(port string, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()

	ctx := context.Background()

	address := "exchange" + port[4:5] + ":" + port
	exchange := model.Exchange{Name: "Exchange" + port[4:5]}

	const maxRetries = 5
	const maxBackoff = 16 * time.Second

	for {
		log.Printf("Подключение к %s", address)

		client := live.NewLiveExchangeClient(exchange, address)
		err := client.Connect()
		if err != nil {
			log.Printf("Ошибка подключения к %s: %v", address, err)
			if !tryReconnect(ctx, client, address, maxRetries, maxBackoff) {
				log.Printf("Не удалось переподключиться к %s после %d попыток", address, maxRetries)
				return
			}
		}

		priceCh := make(chan model.Price, 100)

		// Можно создать отдельный контекст с таймаутом только для чтения (если нужно)
		readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		client.StartReading(readCtx, priceCh)

		for price := range priceCh {
			select {
			case output <- price:
			case <-readCtx.Done():
				log.Printf("Контекст чтения отменен для %s", address)
				break
			}
		}
		cancel() // Закрываем контекст чтения

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
