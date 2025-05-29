package concurrency

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"marketfuck/internal/adapter/in/exchange/live"
	"marketfuck/internal/domain/model"
)

func GenAggr(counter *atomic.Uint64) {
	ports := []string{"40101", "40102", "40103"}

	var wg sync.WaitGroup
	var exchangeWg sync.WaitGroup // Отдельный WaitGroup для бирж
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
	for price := range fanInChan {
		fmt.Println("Получена цена:", price)
	}

	fmt.Println("Программа завершена.")
}
