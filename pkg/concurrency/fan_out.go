package concurrency

import (
	"log"

	"marketfuck/internal/domain/model"
)

func FanOut(input <-chan model.Price, workerChans []chan model.Price) {
	i := 0
	for price := range input {
		workerChans[i%len(workerChans)] <- price
		i++
	}
	log.Println("FanOut: все данные распределены")
}
