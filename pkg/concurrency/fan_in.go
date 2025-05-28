package concurrency



import (
	"marketfuck/internal/domain/model"
	"sync"
)

func FanIn(workerCount int, result chan model.Price, outChannels []chan model.Price, wg *sync.WaitGroup) {
	go func() {
		// Считываем данные из каждого канала воркера
		for i := 0; i < workerCount; i++ {
			wg.Add(1) //?
			go func(i int) {
				defer wg.Done() // Уменьшаем счетчик ожидания, когда горутина завершится
				for price := range outChannels[i] {
					result <- price // Отправляем данные в общий канал
				}
			}(i)
		}
	}()
}
