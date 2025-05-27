package concurrency

// гопота

import (
	"marketfuck/internal/domain/model"
)

func FanIn(workerCount int, input chan model.Price, result chan model.Price) {
	go func() {
		for i := 0; i < workerCount; i++ {
			go func() {
				for price := range input {
					result <- price
				}
			}()
		}
		close(result)
	}()
}
