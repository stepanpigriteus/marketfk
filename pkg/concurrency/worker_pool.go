package concurrency

import (
	"log"
	"marketfuck/internal/domain/model"
	"sync"
	"sync/atomic"
)

type WorkerStr struct {
	In      <-chan model.Price
	Out     chan model.Price
	Id      int
	Counter *atomic.Uint64
	wg      *sync.WaitGroup
}

func NewWorker(id int, in <-chan model.Price, wg *sync.WaitGroup, counter *atomic.Uint64, out chan model.Price) *WorkerStr {
	return &WorkerStr{
		In:      in,
		Out:     out,
		Id:      id,
		Counter: counter,
		wg:      wg,
	}
}

func (w *WorkerStr) Run() {
	defer w.wg.Done()

	for price := range w.In {
		w.Counter.Add(1)
		log.Printf("Worker %d обработал цену: %+v, %d", w.Id, price, w.Counter.Load())
		// Здесь должна быть ебучая логика отправки в канал:

	}
	log.Printf("Worker %d завершает работу, канал закрыт", w.Id)
}

// func Worker(id int, workerCh <-chan model.Price, wg *sync.WaitGroup, counter *atomic.Uint64) {
// 	defer wg.Done()

// 	for price := range workerCh {
// 		counter.Add(1)
// 		log.Printf("Worker %d обработал цену: %+v, %d", id, price, counter.Load())
// 		// Здесь должна быть ебучая логика отправки в канал:

// 	}
// 	log.Printf("Worker %d завершает работу, канал закрыт", id)
// }
