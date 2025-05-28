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
	defer close(w.Out) // grok
	for price := range w.In {
		// w.Counter.Add(1)
		// Здесь должна быть ебучая логика отправки в канал:
		log.Printf("Worker %d: обработана цена %+v", w.Id, price)
		w.Out <- price

	}
	log.Printf("Worker %d завершает работу, канал закрыт", w.Id)
}
