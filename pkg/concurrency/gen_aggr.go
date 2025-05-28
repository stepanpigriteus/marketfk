package concurrency

import (
	"fmt"
	"marketfuck/internal/adapter/in/exchange/live"
	"marketfuck/internal/domain/model"
	"sync"
	"sync/atomic"
)

func GenAggr(counter *atomic.Uint64) {
	ports := []string{"40101", "40102", "40103"}

	var wg sync.WaitGroup
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

	// Запуск Fan-Out для распределения данных по воркерам
	for _, chann := range priceChannels {
		go FanOut(chann, workerChans)
	}

	// Подключение к биржам и получение данных
	for i, port := range ports {
		wg.Add(1)
		go live.GenConnectAndRead(port, &wg, priceChannels[i])
	}

	fanInChan := make(chan model.Price, workerCount*300)
	FanIn(workerCount, fanInChan, outChannels[:], &wg)

	// Ожидаем завершения всех горутин
	go func() {
		wg.Wait()
		close(priceCh1)
		close(priceCh2)
		close(priceCh3)
		close(fanInChan)
		fmt.Println("Все операции завершены.")
	}()

	// Ожидаем, пока данные не будут обработаны
	for price := range fanInChan {
		fmt.Println("Получена цена:", price)
	}

	// Завершаем программу после того, как все данные обработаны и каналы закрыты
	fmt.Println("Программа завершена.")
}
