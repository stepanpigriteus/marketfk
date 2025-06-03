package concurrency

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"marketfuck/internal/adapter/in/exchange/live"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/domain/model"
	"os"
	"sync"
	"sync/atomic"
)

func GenAggr(counter *atomic.Uint64, redis redis.RedisCache) {
	ports := []string{"40101", "40102", "40103"}

	var wg sync.WaitGroup
	var exchangeWg sync.WaitGroup
	priceCh1 := make(chan model.Price, 500)
	priceCh2 := make(chan model.Price, 500)
	priceCh3 := make(chan model.Price, 500)
	priceChannels := [3]chan model.Price{priceCh1, priceCh2, priceCh3}

	const workerCount int = 10
	outChannels := [workerCount]chan model.Price{}
	workerChans := make([]chan model.Price, workerCount)

	for i := 0; i < workerCount; i++ {
		workerChans[i] = make(chan model.Price, 100)
		outChannels[i] = make(chan model.Price, 100)
	}

	// Запуск воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		worker := NewWorker(i, workerChans[i], &wg, counter, outChannels[i])
		go worker.Run()
	}

	// Запуск FanOut с отслеживанием через WaitGroup
	var fanOutWg sync.WaitGroup
	for _, chann := range priceChannels {
		fanOutWg.Add(1)
		go func(chann chan model.Price) {
			FanOut(chann, workerChans)
			fanOutWg.Done()
		}(chann)
	}

	// Закрытие workerChans после завершения всех FanOut
	go func() {
		fanOutWg.Wait()
		for _, ch := range workerChans {
			close(ch)
		}
		log.Println("Все workerChans закрыты после завершения FanOut")
	}()

	// Подключение к биржам с отдельным WaitGroup
	for i, port := range ports {
		exchangeWg.Add(1)
		go func(i int, port string) {
			live.GenConnectAndRead(port, &exchangeWg, priceChannels[i])
		}(i, port)
	}

	// Закрытие priceCh после завершения чтения с бирж
	go func() {
		exchangeWg.Wait()
		close(priceCh1)
		close(priceCh2)
		close(priceCh3)
		log.Println("Закрыты каналы priceCh после завершения чтения с бирж")
	}()

	fanInChan := make(chan model.Price, workerCount*300)
	FanIn(workerCount, fanInChan, outChannels[:], &wg)

	// Закрытие fanInChan после завершения воркеров и FanIn
	go func() {
		wg.Wait()
		close(fanInChan)
		fmt.Println("Все операции завершены.")
	}()

	// Обработка данных

	file, err := os.OpenFile("prices.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// нужно вовзращать значение из функции - работу с редисом вынести прочь
	for price := range fanInChan {

		// Преобразуем price в JSON
		priceJSON, err := json.Marshal(price)
		if err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			continue
		}

		// Создаем ключ
		key := fmt.Sprintf("%s:%s", price.PairName, price.Exchange) // например ETHUSDT:Exchange2

		// Устанавливаем в Redis без TTL (0)
		err = redis.Set(context.Background(), key, string(priceJSON), 0)
		if err != nil {
			log.Printf("Ошибка установки в Redis: %v", err)
		}
		ctx := context.Background()
		keys := "ETHUSDT:Exchange2"

		value, err := redis.Get(ctx, keys)
		if err != nil {
			log.Printf("Ошибка при получении из Redis: %v", err)
		} else {
			var price model.Price
			err = json.Unmarshal([]byte(value), &price)
			if err != nil {
				log.Printf("Ошибка при разборе JSON: %v", err)
			} else {
				fmt.Printf("Извлечено из кеша: %+v\n", price)
			}
		}
	}

	fmt.Println("Программа завершена.")
}
