package concurrency

// гопота

import (
	"log"
	"marketfuck/internal/domain/model"
	"sync"
)

func FanIn(workerCount int, out chan model.Price, inputs []chan model.Price, wg *sync.WaitGroup) {
	var fanInWg sync.WaitGroup

	// Для каждого входного канала запускаем горутину
	for i, inputChan := range inputs {
		fanInWg.Add(1)
		go func(ch <-chan model.Price, workerID int) {
			defer fanInWg.Done()
			for price := range ch {
				out <- price
			}
			log.Printf("FanIn worker %d finished", workerID)
		}(inputChan, i)
	}

	// Ждем завершения всех FanIn горутин и закрываем выходной канал
	go func() {
		fanInWg.Wait()
		close(out)
		log.Println("FanIn completely finished, output channel closed")
	}()
}
